/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 24-08-2018
 * |
 * | File Name:     receiver.go
 * +===============================================
 */

package receivers

// Receiver provides a way to receive data
// when generating load. This option can provide
// ways to better measure performance.
type Receiver struct {
	Topic   []byte
	Handler func(topicName, message []byte)
}
