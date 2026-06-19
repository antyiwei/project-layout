package event

type Handler func(payload any) error

func Publish(topic string, payload any) error {
	// TODO: wire to MQ or in-process bus
	return nil
}

func Subscribe(topic string, handler Handler) {
	// TODO: wire to MQ or in-process bus
}
