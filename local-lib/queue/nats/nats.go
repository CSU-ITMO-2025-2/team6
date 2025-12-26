package queue

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type (
	Nats struct {
		Uri    string `yaml:"uri"`
		Stream string `yaml:"stream"`
		nc     *nats.Conn
		js     nats.JetStreamContext
	}
)

type Queue interface {
	PublishStream(stream string, v any) error
	AddConsumer(stream string)
	AddStream(stream string)
	AddStreamWithSubject(stream string, subjects []string)
	PullSubscribeWithSubject(stream string, subjects string) (*nats.Subscription, error)
	PullSubscribe(stream string) (*nats.Subscription, error)
	QueueSyncSubscribe(subj string) (*nats.Subscription, error)
	Subscribe(stream string, msgHandler func(msg *nats.Msg)) (*nats.Subscription, error)
}

func InitFromEnv() Nats {
	return Nats{
		Uri:    os.Getenv("NATS_URI"),
		Stream: os.Getenv("NATS_STREAM"),
	}
}

func Connect(n Nats) *Nats {
	nc, err := nats.Connect(n.Uri, nats.ClosedHandler(func(c *nats.Conn) {
		log.Fatalf("Exiting: %v", c.LastError())
	}))
	if err != nil {
		log.Fatal("Crit", err)
	}

	n.nc = nc

	n.js, err = nc.JetStream()
	if err != nil {
		log.Fatal("Crit", err)
	}

	return &n
}

func (n *Nats) PublishStream(stream string, v any) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = n.js.Publish(stream, payload)

	return err
}

func (n *Nats) AddStream(stream string) {
	n.js.AddStream(&nats.StreamConfig{
		Name:       stream,
		Subjects:   []string{stream},
		Duplicates: 1 * time.Second,
		Storage:    nats.FileStorage,
		MaxAge:     8 * time.Hour,
		Discard:    nats.DiscardOld,
	})
}

func (n *Nats) AddStreamWithSubject(stream string, subjects []string) {
	subj := []string{stream}

	for _, s := range subjects {
		subj = append(subj, stream+"."+s)
	}

	n.js.AddStream(&nats.StreamConfig{
		Name:       stream,
		Subjects:   subj,
		Duplicates: time.Second,
		Storage:    nats.FileStorage,
		Discard:    nats.DiscardOld,
	})
}

func (n *Nats) AddConsumer(stream string) {
	n.js.AddConsumer(stream, &nats.ConsumerConfig{
		Durable:       stream,
		AckPolicy:     nats.AckAllPolicy,
		DeliverPolicy: nats.DeliverAllPolicy,
	})
}

func (n *Nats) PullSubscribe(stream string) (*nats.Subscription, error) {
	return n.js.PullSubscribe(stream, stream)
}

func (n *Nats) PullSubscribeWithSubject(stream string, subjects string) (*nats.Subscription, error) {
	return n.js.PullSubscribe(stream, subjects)
}

func (n *Nats) Subscribe(stream string, msgHandler func(msg *nats.Msg)) (*nats.Subscription, error) {
	return n.js.Subscribe(stream, msgHandler)
}

func (n *Nats) QueueSyncSubscribe(stream string) (*nats.Subscription, error) {
	return n.nc.QueueSubscribeSync(stream, stream)
}
