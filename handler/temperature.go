package handler

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

func ReceiveTemp(client MQTT.Client, message MQTT.Message) {

}

func HttpTemp(c *fiber.Ctx) (err error) {
	return
}
