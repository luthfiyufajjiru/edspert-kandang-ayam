package handler

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

func ReceiveHumidity(client MQTT.Client, message MQTT.Message) {

}

func HttpHumidity(c *fiber.Ctx) (err error) {
	return
}
