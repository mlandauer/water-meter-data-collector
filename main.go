package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/JuulLabs-OSS/ble"
	"github.com/JuulLabs-OSS/ble/examples/lib/dev"
)

func filter(a ble.Advertisement) bool {
	return a.LocalName() == "Water Meter"
}

func handler(req []byte) {
	fmt.Println(binary.LittleEndian.Uint16(req))
}

var automationIOServiceUUID = ble.UUID16(0x1815)
var analogCharacteristicUUID = ble.UUID16(0x2a58)

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
	services, err := cl.DiscoverServices([]ble.UUID{automationIOServiceUUID})
	if err != nil {
		log.Fatal(err)
	}
	if len(services) == 0 {
		log.Fatal("Couldn't find automation IO service")
	}
	service := services[0]
	characteristics, err := cl.DiscoverCharacteristics([]ble.UUID{analogCharacteristicUUID}, service)
	if err != nil {
		log.Fatal(err)
	}
	if len(characteristics) == 0 {
		log.Fatal("Couldn't find analog characteristic")
	}
	analog := characteristics[0]
	// Looks like (at least on Linux) we also need to explicitly discover the descriptors
	// otherwise later operations fail
	_, err = cl.DiscoverDescriptors(nil, analog)
	if err != nil {
		log.Fatal(err)
	}
	err = cl.Subscribe(analog, false, handler)
	if err != nil {
		log.Fatal(err)
	}
	select {}
}
