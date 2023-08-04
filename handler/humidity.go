package handler

import (
	"bufio"
	"edspert-kandang-ayam/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var strgHum storage.SerialFloat

func ReceiveHumidity(client MQTT.Client, message MQTT.Message) {
	log.Printf("received humidity message with id %d: %s", message.MessageID(), string(message.Payload()))
	strgHum.Append(string(message.Payload()))
}

func humidityEvent(c *fiber.Ctx) (err error) {
	ctx := c.Context()

	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
	event:
		for {
			ln := strgHum.Len()
			data := strgHum.Range(0, ln)
			res := map[string]interface{}{
				"success": true,
				"data":    data,
			}
			msg, _ := json.Marshal(res)
			fmt.Fprintf(w, "data: %s\n\n", string(msg))
			fmt.Println(string(msg))
			err := w.Flush()
			if err != nil {
				log.Println("connection closed")
				break event
			}
			time.Sleep(1 * time.Second)
		}
	}))
	return
}

func httpHumidity(c *fiber.Ctx) (err error) {
	treshold := c.QueryInt("len")
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

func HttpHumidity(c *fiber.Ctx) (err error) {
	switch c.Params("mode") {
	case "event":
		err = humidityEvent(c)
		return
	case "range":
		err = httpHumidity(c)
		return
	default:
		c.Status(http.StatusNotFound).SendString("endpoint not found")
		return
	}
}
