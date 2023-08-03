package handler

import (
	"edspert-kandang-ayam/storage"
	"log"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

var strgTemp storage.SerialFloat

func ReceiveTemp(client MQTT.Client, message MQTT.Message) {
	log.Printf("received temperature message with id %d: %s", message.MessageID(), string(message.Payload()))
	strgTemp.Append(string(message.Payload()))
}

func HttpTemp(c *fiber.Ctx) (err error) {
	treshold := 10
	ln := strgTemp.Len()
	if treshold > ln {
		treshold = ln
	}
	data := strgTemp.Range(ln-treshold, ln)
	res := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	c.Status(http.StatusOK).JSON(res)
	return
}
