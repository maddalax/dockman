package internal

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"log/slog"
	"paas/internal/util"
	"time"
)

type DistributedLock struct {
	key     string
	timeout time.Duration
	c       *KvClient
}

func (c *KvClient) NewLock(key string, timeout time.Duration) *DistributedLock {
	return &DistributedLock{
		key:     key,
		c:       c,
		timeout: timeout,
	}
}

func (l *DistributedLock) Bucket() string {
	return fmt.Sprintf("locks-%s", l.key)
}

func (l *DistributedLock) Lock() error {
	bucket, err := l.c.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: l.Bucket(),
		TTL:    l.timeout,
	})
	if err != nil {
		return err
	}
	_, err = bucket.Create(l.key, []byte("locked"))
	if err != nil {
		if errors.Is(err, nats.ErrKeyExists) {
			// wait for the builderRegistryLock to be released
			success := util.WaitFor(l.timeout, 25*time.Millisecond, func() bool {
				_, err = bucket.Create(l.key, []byte("locked"))
				return err == nil
			})
			if !success {
				return errors.New("builderRegistryLock timeout")
			}
		} else {
			return err
		}
	}
	return nil
}

func (l *DistributedLock) Unlock() error {
	bucket, err := l.c.GetBucket(l.Bucket())
	if err != nil {
		return err
	}
	slog.Debug("unlocking %s", slog.String("key", l.key))
	return bucket.Delete(l.key)
}
