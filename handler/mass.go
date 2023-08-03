package handler

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

func ReceiveMass(client MQTT.Client, message MQTT.Message) {

}

func HttpMass(c *fiber.Ctx) (err error) {
	return
}
