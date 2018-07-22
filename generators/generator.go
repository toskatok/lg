/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-07-2018
 * |
 * | File Name:     generator.go
 * +===============================================
 */

package generators

import "time"

// Generator generates data whenever you want
// based on given input. input can be nil.
type Generator interface {
	Generate(input interface{}) ([]byte, error)
	Topic() []byte
}

// Run given generator on given duration
func Run(g Generator, d time.Duration) {
	sendTick := time.Tick(d)

	for {
		select {
		case <-sendTick:
		}
	}
}
