package cmd

import (
	"log"
	"time"

	"github.com/w6ipa/go-elecraft/rig"
)

func buffRead(k *rig.Connection, dataChan chan []byte, done chan struct{}) {
	cmd := rig.NewTBX()
	//empty buffer
	k.SendCommand(cmd, nil)

	ticker := time.NewTicker(150 * time.Millisecond)
	for {

		select {
		case <-ticker.C:
			data, err := k.SendCommand(cmd, nil)
			if err != nil {
				log.Printf(err.Error())
				close(dataChan)
				return
			}
			if len(data) > 0 {
				dataChan <- data
			}
		}
	}
}
