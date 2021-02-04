package main

import (
	"context"
	"time"

	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cqueue"
)

func NewLogTask(logger clogger.Logger) cqueue.WorkerResult {
	return cqueue.WorkerResult{
		Worker: cqueue.Worker{
			TaskType: "log",
			Timeout:  30 * time.Second,
			Handler: func(ctx context.Context, payload []byte) ([]byte, error) {
				logger.Info("Hello, World!")

				return nil, nil
			},
		},
	}
}
