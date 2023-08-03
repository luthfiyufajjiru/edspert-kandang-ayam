package handler

import (
	"edspert-kandang-ayam/storage"
	"log"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

var strgAmmonia storage.SerialFloat

func ReceiveAmmonia(client MQTT.Client, message MQTT.Message) {
	log.Printf("received ammonia message with id %d: %s", message.MessageID(), string(message.Payload()))
	strgAmmonia.Append(string(message.Payload()))
}

func HttpAmmonia(c *fiber.Ctx) (err error) {
	treshold := 10
	ln := strgAmmonia.Len()
	if treshold > ln {
		treshold = ln
	}
	data := strgAmmonia.Range(ln-treshold, ln)
	res := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	c.Status(http.StatusOK).JSON(res)
	return
}
