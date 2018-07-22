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
}

// Run given generator on given duration
func Run(g Generator, d time.Duration) {
	sendTick := time.Tick(time.Duration(*rate) * time.Millisecond)
	printTick := time.Tick(1 * time.Second)

	// Packet counter
	packets := 0

	for {
		select {
		case <-sendTick:
		case <-printTick:
		}
	}
}
