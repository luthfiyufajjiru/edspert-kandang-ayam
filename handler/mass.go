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

var strgMass storage.SerialFloat

func ReceiveMass(client MQTT.Client, message MQTT.Message) {
	log.Printf("received mass message with id %d: %s", message.MessageID(), string(message.Payload()))
	strgMass.Append(string(message.Payload()))
}

func massEvent(c *fiber.Ctx) (err error) {
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
			ln := strgMass.Len()
			data := strgMass.Range(0, ln)
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

func httpMass(c *fiber.Ctx) (err error) {
	treshold := c.QueryInt("len")
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

func HttpMass(c *fiber.Ctx) (err error) {
	switch c.Params("mode") {
	case "event":
		err = massEvent(c)
		return
	case "range":
		err = httpMass(c)
		return
	default:
		c.Status(http.StatusNotFound).SendString("endpoint not found")
		return
	}
}
