package queue

import "github.com/nats-io/nats.go"

const (
	streamMLService = "core-task-ml"
	requestSubject  = "request"
	responseSubject = "response"
)

type Study interface {
	Queue
	PullSubscribeML() (*nats.Subscription, error)
	PublishStreamML(msg any) error
	AddMLStream()
}

type Wrapper struct {
	*Nats
}

func WrapNatsCore(nats *Nats) *Wrapper {
	return &Wrapper{nats}
}

func (n *Wrapper) PullSubscribeML() (*nats.Subscription, error) {
	n.AddStreamWithSubject(streamMLService, []string{requestSubject, responseSubject})
	n.AddConsumer(streamMLService)
	return n.PullSubscribeWithSubject(streamMLService, requestSubject)
}
func (n *Wrapper) PublishStreamML(msg any) error {
	return n.PublishStream(streamMLService+"."+requestSubject, msg)
}

func (n *Wrapper) AddMLStream() {
	n.AddStreamWithSubject(streamMLService, []string{requestSubject, responseSubject})
}
