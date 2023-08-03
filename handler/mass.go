package handler

import (
	"edspert-kandang-ayam/storage"
	"log"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

var strgMass storage.SerialFloat

func ReceiveMass(client MQTT.Client, message MQTT.Message) {
	log.Printf("received mass message with id %d: %s", message.MessageID(), string(message.Payload()))
	strgMass.Append(string(message.Payload()))
}

func HttpMass(c *fiber.Ctx) (err error) {
	treshold := 10
	ln := strgMass.Len()
	if treshold > ln {
		treshold = ln
	}
	data := strgMass.Range(ln-treshold, ln)
	res := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	c.Status(http.StatusOK).JSON(res)
	return
}
