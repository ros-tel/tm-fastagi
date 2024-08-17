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
		js     nats.JetStreamContext
		kv     nats.KeyValue
	}
)

func Connect(n Nats) *Nats {
	nc, err := nats.Connect(n.Uri)
	if err != nil {
		log.Fatal("Crit", err)
		os.Exit(1)
	}

	n.js, err = nc.JetStream()
	if err != nil {
		log.Fatal("Crit", err)
		os.Exit(1)
	}

	return &n
}

func (n *Nats) PublishStream(stream string, v interface{}) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = n.js.Publish(stream, payload)

	return err
}

func (n *Nats) AddStream(stream string) {
	_, err := n.js.AddStream(&nats.StreamConfig{
		Name:       stream,
		Subjects:   []string{stream},
		Duplicates: time.Duration(1 * time.Second),
		Storage:    nats.FileStorage,
		MaxAge:     time.Duration(8 * time.Hour),
		Discard:    nats.DiscardOld,
	})
	if err != nil {
		log.Fatal("Crit", err)
		os.Exit(1)
	}
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
