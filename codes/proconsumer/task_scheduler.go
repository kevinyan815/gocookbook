package proconsumer

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type TaskScheduler struct {
	errGroup *errgroup.Group
	producer *Producer
	consumer *Consumer
	consumerNum int
}

func NewTaskScheduler (pFunc ProducerCallerFunc, cFunc ConsumerCallerFunc, consumerNum int) (ts *TaskScheduler) {

	eg, cancelCtx := errgroup.WithContext(context.Background())

	dataCh := make(chan interface{}, 1000)
	producer := &Producer{
		DataCh:    dataCh,
	}
	consumer := &Consumer{
		DataCh: dataCh,
	}
	producer.RegisterProduceFunc(pFunc)
	consumer.RegisterConsumeFunc(cFunc)
	pCancelCtx, _ := context.WithCancel(cancelCtx)
	cCancelCtx, _ := context.WithCancel(cancelCtx)
	producer.CancelCtx = pCancelCtx
	consumer.CancelCtx = cCancelCtx
	ts = &TaskScheduler{
		errGroup: eg,
		producer: producer,
		consumer: consumer,
		consumerNum: consumerNum,
	}

	return
}

func (ts *TaskScheduler) Execute () (err error) {

	ts.errGroup.Go(ts.producer.Run)

	for i := 0; i < ts.consumerNum; i++ {
		ts.errGroup.Go(ts.consumer.Run)
	}

	err = ts.errGroup.Wait()
	if err != nil {
		// 如果有goroutine执行发生错误, 等待2秒尽可能保证其他goroutine完成收尾工作。
		time.Sleep(time.Second * 2)
	}
	return err
}
