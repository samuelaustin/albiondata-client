package client

import (
	"github.com/samuelaustin/albiondata-client/lib"
	"github.com/samuelaustin/albiondata-client/log"
	//"github.com/samuelaustin/albiondata-client/notification"
	"encoding/json"
	"fmt"
	"time"
	"io/ioutil"
)

type eventBankContainerContents struct {
	Slots       []int64         `mapstructure:"1"`
	ContainerID lib.CharacterID `mapstructure:"2"`
}

func (event eventBankContainerContents) Process(state *albionState) {
	log.Debug("Got Bank container opening event...")

	items := []*lib.ItemContainer{}

	for _, v := range event.Slots {
		if v == 0 {
			continue
		}

		item := state.ContainerItemsToSend[v]
		items = append(items, &item)
	}

	// clear the array to prepare for the next container opening
	state.ContainerItemsToSend = make(map[int64]lib.ItemContainer)

	if len(items) < 1 {
		return
	}

	/*if state.LocationId == 0 {
		msg := "The players location has not yet been set. Please transition zones so the location can be identified."
		log.Warn(msg)
		notification.Push(msg)
		return
	}*/

	/*upload := lib.ContainerUpload{
		Items:           items,
		CurrentLocation: state.LocationId,
		ContainerType:   "Bank",
		ContainerGUID:   event.ContainerID,
	}*/

	//log.Infof("Sending Bank container with %d items of %v to ingest.", len(items), state.CharacterName)

	log.Infof("Sending contents of chest to json file...")

	itemsJson, _ := json.Marshal(items)
	date := time.Now().Local().Format("2006-01-02T15:04:05")
	filename := fmt.Sprintf("%v_%v.json", event.ContainerID, date)
    ioutil.WriteFile(filename, itemsJson, 0644)


	//sendMsgToPrivateUploaders(&upload, lib.NatsBankContainerData, state)
}