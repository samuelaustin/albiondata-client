package client

import (
	"runtime"
	"time"

	"github.com/samuelaustin/albiondata-client/log"
)

type processesWatcher struct {
	knownAlbions   []int
	albionWatchers map[int]*albionProcessWatcher
}

func newProcessWatcher() *processesWatcher {
	return &processesWatcher{
		albionWatchers: make(map[int]*albionProcessWatcher),
	}
}

func (pw *processesWatcher) run() {
	log.Debug("Watching processes for Albion to start...")

	for {
		var process_string string

		if runtime.GOOS == "windows" {
			process_string = "Albion-Online.exe"
		} else {
			process_string = "Albion-Online"
		}

		current := findProcess(process_string)

		added, removed := diffIntSets(pw.knownAlbions, current)

		for _, pid := range added {
			apw := newAlbionProcessWatcher(pid)

			pw.albionWatchers[pid] = apw

			go apw.run()
		}

		for _, pid := range removed {
			pw.albionWatchers[pid].quit <- true
			delete(pw.albionWatchers, pid)
		}

		pw.knownAlbions = current
		time.Sleep(time.Second)
	}
}
