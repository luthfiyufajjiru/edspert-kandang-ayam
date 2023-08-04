package main

import (
	"edspert-kandang-ayam/handler"
	"edspert-kandang-ayam/mqtt"
	"fmt"
	"log"
	"os"
	"os/signal"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

func main() {
	server := "broker.hivemq.com"
	port := 1883
	clientID := "go-client"

	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", server, port))
	opts.SetClientID(clientID)

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
	}
	defer client.Disconnect(250)

	fmt.Printf("Connected to %s:%d\n", server, port)

	// Start the MQTT subscriber logic in a separate Goroutine
	go mqtt.Subscriber(client, "edspertkel3/amo", handler.ReceiveAmmonia)
	go mqtt.Subscriber(client, "edspertkel3/mass", handler.ReceiveMass)
	go mqtt.Subscriber(client, "edspertkel3/hum", handler.ReceiveHumidity)
	go mqtt.Subscriber(client, "edspertkel3/temp", handler.ReceiveTemp)

	// Start the Fiber HTTP server
	app := fiber.New()

	app.Static("", "./public")
	v1 := app.Group("/api/v1")

	v1.Get("/ammonia/:mode", handler.HttpAmmonia)
	v1.Get("/mass/:mode", handler.HttpMass)
	v1.Get("/temperature/:mode", handler.HttpTemp)
	v1.Get("/humidity/:mode", handler.HttpHumidity)

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	// Wait for an interruption signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Exited")
}
