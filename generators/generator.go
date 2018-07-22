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

// Generator generates data whenever you want
// based on given input. input can be nil.
type Generator interface {
	Generate(input interface{}) ([]byte, error)
	Topic() []byte
}
