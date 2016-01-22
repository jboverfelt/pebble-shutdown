package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/mqtt"
	"github.com/hybridgroup/gobot/platforms/pebble"
)

const (
	UP     = "up"
	DOWN   = "down"
	SELECT = "select"
)

type curCmds struct {
	mu   *sync.Mutex
	cmds []string
}

func clearCmds(c *curCmds) {
	for {
		<-time.After(5 * time.Second)
		c.mu.Lock()
		c.cmds = make([]string, 5)
		c.mu.Unlock()
	}
}

func isCombo(c *curCmds, btn string) bool {
	c.cmds = append(c.cmds, btn)
	l := len(c.cmds)
	if btn == SELECT && l >= 5 {
		fourth := c.cmds[l-2]
		third := c.cmds[l-3]
		second := c.cmds[l-4]
		first := c.cmds[l-5]

		return fourth == DOWN && third == DOWN && second == UP && first == UP
	}

	return false
}

func pebbleWork(pebbleDriver *pebble.PebbleDriver, mqttAdaptor *mqtt.MqttAdaptor, c *curCmds) func() {
	return func() {
		gobot.On(pebbleDriver.Event("button"), func(data interface{}) {
			btn := data.(string)
			c.mu.Lock()
			defer c.mu.Unlock()

			if isCombo(c, btn) {
				ok := mqttAdaptor.Publish("chip/command/shutdown", []byte("shutdown"))
				if !ok {
					panic("Error publishing message")
				}

				log.Printf("Received combo, enqueued shutdown message")
			}
		})
	}
}

func main() {
	host := flag.String("host", "0.0.0.0:1883", "Hostname and port of the MQTT Broker")
	flag.Parse()

	gbot := gobot.NewGobot()
	api.NewAPI(gbot).Start()

	pebbleAdaptor := pebble.NewPebbleAdaptor("pebble")
	pebbleDriver := pebble.NewPebbleDriver(pebbleAdaptor, "pebble")
	mqttAdaptor := mqtt.NewMqttAdaptor("server", fmt.Sprintf("tcp://%s", *host), "shutdowner")

	c := &curCmds{mu: &sync.Mutex{}, cmds: make([]string, 5)}

	go clearCmds(c)

	pebbleRobot := gobot.NewRobot("pebble",
		[]gobot.Connection{pebbleAdaptor},
		[]gobot.Device{pebbleDriver},
		pebbleWork(pebbleDriver, mqttAdaptor, c),
	)

	mqttRobot := gobot.NewRobot("mqtt",
		[]gobot.Connection{mqttAdaptor},
	)

	gbot.AddRobot(mqttRobot)
	gbot.AddRobot(pebbleRobot)

	errs := gbot.Start()

	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
