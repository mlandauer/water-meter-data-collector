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
	fmt.Printf("distance: %v cm\n", binary.LittleEndian.Uint16(req))
}

func batteryLevelHandler(req []byte) {
	fmt.Printf("battery: %v\n", req[0])
}

var batteryLevelCharacteristicUUID = ble.UUID16(0x2a19)
var analogCharacteristicUUID = ble.UUID16(0x2a58)

func main() {
	device, err := dev.NewDevice("default")
	if err != nil {
		log.Fatal(err)
	}
	ble.SetDefaultDevice(device)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Connecting to water meter...")
	cl, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	done := make(chan struct{})
	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-cl.Disconnected()
		log.Println("Disconnected")
		close(done)
	}()

	profile, err := cl.DiscoverProfile(true)
	if err != nil {
		log.Fatal(err)
	}

	analog := profile.FindCharacteristic(ble.NewCharacteristic(analogCharacteristicUUID))
	if analog == nil {
		log.Fatal("Couldn't find analog characteristic")
	}

	err = cl.Subscribe(analog, false, handler)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Subscribed to analog notifications")

	batteryLevel := profile.FindCharacteristic(ble.NewCharacteristic(batteryLevelCharacteristicUUID))
	if batteryLevel == nil {
		log.Fatal("Couldn't find battery level characteristic")
	}

	err = cl.Subscribe(batteryLevel, false, batteryLevelHandler)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Subscribed to battery notifications")

	<-done
}
