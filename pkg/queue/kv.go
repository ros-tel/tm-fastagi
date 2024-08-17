package queue

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// Создание
func (n *Nats) CreateKeyValue(stream string, ttl time.Duration) {
	var err error
	n.kv, err = n.js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:  stream,
		TTL:     ttl,
		History: 1,
	})
	if err != nil {
		log.Fatal("Crit", err)
		os.Exit(1)
	}
}

// Использование не создавая
func (n *Nats) KeyValue(stream string) {
	var err error
	n.kv, err = n.js.KeyValue(stream)
	if err != nil {
		log.Fatal("Crit", err)
		os.Exit(1)
	}
}

func (n *Nats) KvPutString(key string, value string) {
	_, err := n.kv.PutString(key, value)
	if err != nil {
		log.Fatal("Crit", err)
		os.Exit(1)
	}
}

func (n *Nats) KvGetString(key string) (value string, err error) {
	e, err := n.kv.Get(key)
	if err != nil {
		return
	}
	value = string(e.Value())
	return
}

func (n *Nats) KvDelete(key string) error {
	return n.kv.Purge(key)
}
