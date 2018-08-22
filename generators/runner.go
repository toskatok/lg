/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-07-2018
 * |
 * | File Name:     runner.go
 * +===============================================
 */

package generators

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Source is called in order to create source data for generate
type Source func() interface{}

// Runner runs given generator in specific intervals
type Runner struct {
	generator Generator
	duration  time.Duration
	counter   int64
	source    Source

	cli *client.Client
	lck sync.RWMutex

	stop chan struct{}
}

// NewRunner creates new runner
func NewRunner(g Generator, d time.Duration, s Source, broker string) (Runner, error) {
	var r = Runner{
		generator: g,
		duration:  d,
		counter:   0,
		source:    s,

		stop: make(chan struct{}),
	}

	// Create an MQTT Client.
	r.cli = client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			log.Println(err)
		},
	})

	// Connect to the MQTT Server.
	if err := r.cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  broker,
		ClientID: []byte(fmt.Sprintf("I1820-lg-%d", rand.Intn(1024))),
	}); err != nil {
		return r, err
	}

	return r, nil
}

// Count returns number of generated messages
func (r *Runner) Count() int64 {
	r.lck.RLock()
	defer r.lck.RUnlock()
	return r.counter
}

// Stop stops the runner
func (r *Runner) Stop() {
	close(r.stop)
}

// Run runs the runner :joy:
func (r *Runner) Run() {
	sendTick := time.Tick(r.duration)

	go func() {
		for {
			select {
			case <-sendTick:
				message, err := r.generator.Generate(r.source())
				if err != nil {
					log.Printf("Generator Generate: %s", err)
				}

				if err := r.cli.Publish(&client.PublishOptions{
					QoS:       mqtt.QoS0,
					TopicName: r.generator.Topic(),
					Message:   message,
				}); err != nil {
					log.Printf("MQTT Publish: %s", err)
				}

				r.lck.Lock()
				r.counter++
				r.lck.Unlock()
			case <-r.stop:
				return
			}
		}
	}()
}
