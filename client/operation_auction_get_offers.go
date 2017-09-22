package client

import (
	"encoding/json"

	"github.com/samuelaustin/albiondata-client/lib"
	"github.com/samuelaustin/albiondata-client/log"
	"github.com/samuelaustin/albiondata-client/notification"
)

type operationAuctionGetOffers struct {
	Category         string   `mapstructure:"1"`
	SubCategory      string   `mapstructure:"2"`
	Quality          string   `mapstructure:"3"`
	Enchantment      uint32   `mapstructure:"4"`
	EnchantmentLevel string   `mapstructure:"8"`
	ItemIds          []uint16 `mapstructure:"6"`
	MaxResults       uint32   `mapstructure:"9"`
	IsAscendingOrder bool     `mapstructure:"11"`
}

func (op operationAuctionGetOffers) Process(state *albionState) {
	log.Debug("Got AuctionGetOffers operation...")
}

type operationAuctionGetOffersResponse struct {
	MarketOrders []string `mapstructure:"0"`
}

func (op operationAuctionGetOffersResponse) Process(state *albionState) {
	log.Debug("Got response to AuctionGetOffers operation...")

	if state.LocationId == 0 {
		log.Error("The players location has not yet been set. Please transition zones so the location can be identified.")
		notification.Push("The players location has not yet been set. Please transition zones so the location can be identified.")

		return
	}

	orders := []*lib.MarketOrder{}

	for _, v := range op.MarketOrders {
		order := &lib.MarketOrder{}

		err := json.Unmarshal([]byte(v), order)
		if err != nil {
			log.Errorf("Problem converting market order to internal struct: %v", err)
		}
		order.LocationID = state.LocationId
		orders = append(orders, order)
	}

	if len(orders) < 1 {
		return
	}

	upload := lib.MarketUpload{
		Orders: orders,
	}

	log.Infof("Sending %d market offers to ingest", len(orders))
	sendMsgToPublicUploaders(upload, lib.NatsMarketOrdersIngest, state)
}
