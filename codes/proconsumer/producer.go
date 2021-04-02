package proconsumer

import (
	"context"
)

type Producer struct {
	DataCh    chan <- interface{}
	closure   ProducerCallerFunc
	CancelCtx context.Context
}


func (pr *Producer) RegisterProduceFunc(closure ProducerCallerFunc) error {
	pr.closure = closure
	return nil
}

func (pr *Producer) Run() error {
	err := pr.closure(pr)
	return err
}


// 生产者要调用的方法
type ProducerCallerFunc func(producer *Producer) error