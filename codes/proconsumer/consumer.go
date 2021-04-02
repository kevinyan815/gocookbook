package proconsumer

import (
	"context"
)

type Consumer struct {
	DataCh    <- chan interface{}
	closure   ConsumerCallerFunc
	CancelCtx context.Context
}

func (cs *Consumer) RegisterConsumeFunc(consumer ConsumerCallerFunc) error {
	cs.closure = consumer
	return nil
}

func (cs *Consumer) Run () error {
	err := cs.closure(cs)
	return err
}

type ConsumerCallerFunc func (consumer *Consumer) error