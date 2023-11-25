package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/Wolechacho/ticketmaster-backend/common"
	"github.com/nats-io/nats.go"
)

func main() {
	ch := make(chan bool)
	// connect to nats
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("starting publish test")

	//sample message
	bk := common.SeatAvailableMessage{
		CinemaSeatIds: []string{"210a0000-1d2c-6361-3762-376363662d64", "210a0000-1d2c-6662-3636-366564622d39", "210a0000-1d2c-6662-3636-366564622d44"},
		ShowId:        "210a0000-642c-6361-3162-326437372d30",
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(bk)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	if err := nc.Publish(common.SeatAvailableEvent, buf.Bytes()); err != nil {
		log.Fatal(err)
	}

	nc.Flush()

	<-ch
}
