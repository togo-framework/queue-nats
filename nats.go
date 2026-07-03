// Package queuenats is a togo queue backend backed by NATS. It implements
// togo.Queue over NATS subjects and takes over k.Queue when selected with
// `togo provider:use queue nats` (or TOGO_QUEUE_PROVIDER=nats).
//
//   NATS_URL      NATS server URL (default nats://127.0.0.1:4222)
//   NATS_PREFIX   subject prefix   (default "togo.queue")
package queuenats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/togo-framework/providers"
	"github.com/togo-framework/togo"
)

func init() {
	togo.RegisterProviderFunc("queue-nats", togo.PriorityLate+1, func(k *togo.Kernel) error {
		if providers.Active(k, providers.CapQueue, "memory") != "nats" {
			return nil
		}
		q, err := newNats(k)
		if err != nil {
			if k.Log != nil {
				k.Log.Error("queue-nats: connect failed", "err", err)
			}
			return nil
		}
		k.Queue = q
		if k.Log != nil {
			k.Log.Info("plugin active", "plugin", "queue-nats")
		}
		return nil
	})
}

type natsQueue struct {
	nc     *nats.Conn
	prefix string
}

func newNats(k *togo.Kernel) (*natsQueue, error) {
	url := providers.Value(k, providers.CapQueue, "nats", "url", nats.DefaultURL, false)
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &natsQueue{nc: nc, prefix: providers.Value(k, providers.CapQueue, "nats", "prefix", "togo.queue", false)}, nil
}

func (q *natsQueue) subject(name string) string { return q.prefix + "." + name }

func (q *natsQueue) Handle(name string, h togo.QueueHandler) {
	_, _ = q.nc.Subscribe(q.subject(name), func(m *nats.Msg) {
		var payload any
		_ = json.Unmarshal(m.Data, &payload)
		_ = h(context.Background(), payload)
	})
}

func (q *natsQueue) Dispatch(_ context.Context, name string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return q.nc.Publish(q.subject(name), b)
}
