package transport

import "github.com/nats-io/nats.go"

type NATS struct {
	conn *nats.Conn
}

func (n *NATS) Init(url string, token string) error {
	nc, err := nats.Connect(url)
	if err != nil {
		return err
	}

	n.conn = nc

	return nil
}

func (n *NATS) Transmit(topic string, data []byte) error {
	return n.conn.Publish(topic, data)
}
