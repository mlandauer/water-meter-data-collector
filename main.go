package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/JuulLabs-OSS/ble"
	"github.com/JuulLabs-OSS/ble/examples/lib/dev"
)

func filter(a ble.Advertisement) bool {
	return a.LocalName() == "Water Meter"
}

func handler(req []byte) {
	fmt.Println(binary.LittleEndian.Uint16(req))
}

func main() {
	device, err := dev.NewDevice("default")
	if err != nil {
		log.Fatal(err)
	}
	ble.SetDefaultDevice(device)

	ctx := context.Background()
	cl, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	services, err := cl.DiscoverServices([]ble.UUID{ble.UUID16(0x1815)})
	if err != nil {
		log.Fatal(err)
	}
	service := services[0]
	characteristics, err := cl.DiscoverCharacteristics([]ble.UUID{}, service)
	if err != nil {
		log.Fatal(err)
	}
	// We're only expecting one characteristic - analog
	analog := characteristics[0]
	err = cl.Subscribe(analog, false, handler)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(60 * time.Second)
}
