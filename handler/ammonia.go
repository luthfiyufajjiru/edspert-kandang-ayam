package handler

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

func ReceiveAmmonia(client MQTT.Client, message MQTT.Message) {

}

func HttpAmmonia(c *fiber.Ctx) (err error) {
	return
}
