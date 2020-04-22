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

package runner

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/toskatok/lg/generator"
	"github.com/toskatok/lg/transport"
)

// Pick is called in order to pickup a data for generate method
type Pick func() interface{}

// Transport transports data to given topic based its network protocol
type Transport interface {
	Init(url string, token string) error
	Transmit(topic string, data []byte) error
}

// Config contains runner configuration
// these configuration specifies host, rate and etc.
type Config struct {
	Generator generator.Generator
	Duration  time.Duration
	Pick      Pick
	Token     string
	URL       string
}

// Runner runs given generator in specific intervals
type Runner struct {
	Generator generator.Generator
	Duration  time.Duration
	Pick      Pick

	Transport Transport

	lck     sync.RWMutex
	counter int64

	stop chan struct{}
}

// New creates new runner
func New(config Config) (Runner, error) {
	// Find and configure the transport
	var t Transport

	url, err := url.Parse(config.URL)
	if err != nil {
		return Runner{}, err
	}

	switch url.Scheme {
	case "http", "https":
		t = &transport.HTTP{}
	case "mqtt":
		t = &transport.MQTT{}
	case "kafka":
		t = &transport.Kafka{}
	default:
		return Runner{}, fmt.Errorf("scheme %s is not supported yet", url.Scheme)
	}

	if err := t.Init(url.Host, config.Token); err != nil {
		return Runner{}, err
	}

	return Runner{
		Generator: config.Generator,
		Duration:  config.Duration,
		counter:   0,
		Pick:      config.Pick,

		Transport: t,

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
	ticker := time.NewTicker(r.Duration)

	go func() {
		for {
			select {
			case <-ticker.C:
				message, err := r.Generator.Generate(r.Pick())
				if err != nil {
					log.Printf("Generator Generate: %s", err)
				}

				if err := r.Transport.Transmit(
					r.Generator.Topic(),
					message,
				); err != nil {
					log.Printf("Transmit: %s", err)
				}

				r.lck.Lock()
				r.counter++
				r.lck.Unlock()
			case <-r.stop:
				ticker.Stop()
				return
			}
		}
	}()
}
