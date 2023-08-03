package mqtt

import (
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Subscriber(client MQTT.Client, topic string, callback func(client MQTT.Client, message MQTT.Message)) {
	if token := client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic: %v", token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
