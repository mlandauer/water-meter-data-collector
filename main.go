package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/JuulLabs-OSS/ble"
	"github.com/JuulLabs-OSS/ble/darwin"
)

func filter(a ble.Advertisement) bool {
	return a.LocalName() == "Water Meter"
}

// func handler(a ble.Advertisement) {
// 	fmt.Println("local name:", a.LocalName())
// 	fmt.Println("RSSI:", a.RSSI())
// }

func handler(req []byte) {
	fmt.Println(binary.LittleEndian.Uint16(req))
}

func main() {
	// First let's just try doing some scanning for things that are advertising
	device, err := darwin.NewDevice()
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
