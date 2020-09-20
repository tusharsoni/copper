package cqueue

import "time"

type Config struct {
	DequeueWaitTime time.Duration
}

func (c Config) isValid() bool {
	return int64(c.DequeueWaitTime) > 0
}

func GetDefaultConfig() Config {
	return Config{DequeueWaitTime: 5 * time.Second}
}
