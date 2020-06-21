package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/w6ipa/go-elecraft/rig"
)

func buffRead(k *rig.Connection, dataChan chan []byte, done chan struct{}) {
	// if rig is k3s use TB instead of TBX
	r, err := getRig(k)
	if err != nil {
		log.Print(err)
		return
	}
	var cmd rig.ReqResponder
	if r == "KX3" || r == "KX2" {
		cmd = rig.NewTBX()
	} else {
		log.Printf("Buffered read does not work on K3S")
		close(dataChan)
		return
	}
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

func getRig(k *rig.Connection) (r string, err error) {
	cmd := rig.NewOM()

	buff, err := k.SendCommand(cmd, nil)
	if err != nil {
		return
	}

	rsp := cmd.Parse(buff)
	m, ok := rsp.(map[string]bool)
	if !ok {
		return r, fmt.Errorf("unexpected structure")
	}
	r = rig.Rig(m)

	if len(r) == 0 {
		return r, fmt.Errorf("unknown rig")
	}
	return
}
