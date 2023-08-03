package handler

import (
	"edspert-kandang-ayam/storage"
	"log"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

var strgHum storage.SerialFloat

func ReceiveHumidity(client MQTT.Client, message MQTT.Message) {
	log.Printf("received humidity message with id %d: %s", message.MessageID(), string(message.Payload()))
	strgHum.Append(string(message.Payload()))
}

func HttpHumidity(c *fiber.Ctx) (err error) {
	treshold := 10
	ln := strgHum.Len()
	if treshold > ln {
		treshold = ln
	}
	data := strgHum.Range(ln-treshold, ln)
	res := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	c.Status(http.StatusOK).JSON(res)
	return
}
