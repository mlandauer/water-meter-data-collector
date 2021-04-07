package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JuulLabs-OSS/ble"
	"github.com/JuulLabs-OSS/ble/examples/lib/dev"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	waterHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "water",
		Name:      "water_height_uncalibrated",
		Help:      "Water height in arbitrary units",
	})
	battery = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "water",
		Name:      "battery",
		Help:      "Percentage full of the battery",
	})
)

func filter(a ble.Advertisement) bool {
	return a.LocalName() == "Water Meter"
}

func handler(req []byte) {
	value := binary.LittleEndian.Uint16(req)
	fmt.Printf("water depth (uncalibrated): %v\n", value)
	waterHeight.Set(float64(value))
}

func batteryLevelHandler(req []byte) {
	value := req[0]
	fmt.Printf("battery: %v\n", value)
	battery.Set(float64(value))
}

var batteryLevelCharacteristicUUID = ble.UUID16(0x2a19)
var analogCharacteristicUUID = ble.UUID16(0x2a58)

func captureAndRecord() {
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

func main() {
	go captureAndRecord()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
