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

package core

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/I1820/lg/generators"
	"github.com/I1820/lg/receivers"
)

// Source is called in order to create source data for generate
type Source func() interface{}

// Transport transports data to given topic based its network protocol
type Transport interface {
	Init(url string) error
	Transmit(topic string, data []byte) error
}

// Runner runs given generator in specific intervals
type Runner struct {
	generator generators.Generator
	duration  time.Duration
	counter   int64
	source    Source

	transport Transport

	lck sync.RWMutex

	stop chan struct{}
}

// NewRunner creates new runner
func NewRunner(g generators.Generator, d time.Duration, s Source, rawurl string, rs ...receivers.Receiver) (Runner, error) {
	// TODO correct receivers
	/*
		for _, r := range rs {
			if err := cli.Subscribe(&client.SubscribeOptions{
				SubReqs: []*client.SubReq{
					&client.SubReq{
						TopicFilter: r.Topic,
						QoS:         mqtt.QoS0,
						Handler:     r.Handler,
					},
				},
			}); err != nil {
				return Runner{}, err
			}
		}
	*/

	// Find and configure the transport
	var t Transport
	url, err := url.Parse(rawurl)
	if err != nil {
		return Runner{}, err
	}
	switch url.Scheme {
	case "http", "https":
		t = &HTTPTransport{}
	case "mqtt":
		t = &MQTTTransport{}
	}
	if err := t.Init(rawurl); err != nil {
		return Runner{}, err
	}

	return Runner{
		generator: g,
		duration:  d,
		counter:   0,
		source:    s,

		transport: t,

		stop: make(chan struct{}),
	}, nil
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

				if err := r.transport.Transmit(
					r.generator.Topic(),
					message,
				); err != nil {
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
