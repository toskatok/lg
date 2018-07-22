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
	"log"
	"sync"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Source is called in order to create source data for generate
type Source func() interface{}

// Runner runs given generator in specific intervals
type Runner struct {
	g Generator
	d time.Duration
	p int64
	s Source

	cli *client.Client
	lck sync.RWMutex
}

// NewRunner creates new runner
func NewRunner(g Generator, d time.Duration, s Source, cli *client.Client) Runner {
	return Runner{
		g: g,
		d: d,
		p: 0,
		s: s,

		cli: cli,
	}
}

// Count returns number of generated messages
func (r *Runner) Count() int64 {
	r.lck.RLock()
	defer r.lck.RUnlock()
	return r.p
}

// Run runs the runner :joy:
func (r *Runner) Run() {
	sendTick := time.Tick(r.d)

	go func() {
		for {
			select {
			case <-sendTick:
				message, err := r.g.Generate(r.s())
				if err != nil {
					log.Printf("Generator Generate: %s", err)
				}

				if err := r.cli.Publish(&client.PublishOptions{
					QoS:       mqtt.QoS0,
					TopicName: r.g.Topic(),
					Message:   message,
				}); err != nil {
					log.Printf("MQTT Publish: %s", err)
				}

				r.lck.Lock()
				r.p++
				r.lck.Unlock()
			}
		}
	}()
}
