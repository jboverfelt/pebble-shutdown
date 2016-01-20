package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/pebble"
)

const (
	UP               = "up"
	DOWN             = "down"
	SELECT           = "select"
	SHUTDOWN_MESSAGE = "Shutdown command received"
)

type curCmds struct {
	mu   *sync.Mutex
	cmds []string
}

func clearCmds(c *curCmds, stop <-chan bool) {
	for {
		select {
		case <-time.After(5 * time.Second):
			c.mu.Lock()
			c.cmds = make([]string, 5)
			c.mu.Unlock()
		case <-stop:
			return
		}
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

func main() {
	gbot := gobot.NewGobot()
	api.NewAPI(gbot).Start()

	pebbleAdaptor := pebble.NewPebbleAdaptor("pebble")
	pebbleDriver := pebble.NewPebbleDriver(pebbleAdaptor, "pebble")

	stop := make(chan bool)
	c := &curCmds{mu: &sync.Mutex{}, cmds: make([]string, 5)}

	go clearCmds(c, stop)

	work := func() {
		gobot.On(pebbleDriver.Event("button"), func(data interface{}) {
			btn := data.(string)
			c.mu.Lock()
			defer c.mu.Unlock()

			if isCombo(c, btn) {
				msg := pebbleDriver.SendNotification(SHUTDOWN_MESSAGE)
				log.Printf("Received combo: %s", msg)

				stop <- true

				err := exec.Command("sudo", "shutdown", "-h", "1").Run()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		})
	}

	robot := gobot.NewRobot("pebble",
		[]gobot.Connection{pebbleAdaptor},
		[]gobot.Device{pebbleDriver},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
