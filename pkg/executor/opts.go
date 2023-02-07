// Copyright 2023 NJWS Inc.

package executor

import (
	"time"
)

type Opt func(*executor) error

// WithBroker can set only when creating
func WithBroker(broker string) Opt {
	return func(e *executor) (err error) {
		e.brokers = append(e.brokers, broker)
		return
	}
}

func WithTimeout(timeout time.Duration) Opt {
	return func(e *executor) (err error) {
		e.timeout = timeout
		return
	}
}

func WithTopic(topic string) Opt {
	return func(e *executor) (err error) {
		e.topic = topic
		return
	}
}
