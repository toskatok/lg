package transport

type NATS struct {
}

func (n NATS) Init(url string, token string) error {
	panic("implement me")
}

func (n NATS) Transmit(topic string, data []byte) error {
	panic("implement me")
}
