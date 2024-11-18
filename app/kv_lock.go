package app

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"log/slog"
	"paas/util"
	"time"
)

type Lock struct {
	key     string
	timeout time.Duration
	c       *KvClient
}

func (c *KvClient) NewLock(key string, timeout time.Duration) *Lock {
	return &Lock{
		key:     key,
		c:       c,
		timeout: timeout,
	}
}

func (l *Lock) Bucket() string {
	return fmt.Sprintf("locks-%s", l.key)
}

func (l *Lock) Lock() error {
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
			// wait for the lock to be released
			success := util.WaitFor(l.timeout, 25*time.Millisecond, func() bool {
				_, err = bucket.Create(l.key, []byte("locked"))
				return err == nil
			})
			if !success {
				return errors.New("lock timeout")
			}
		} else {
			return err
		}
	}
	return nil
}

func (l *Lock) Unlock() error {
	bucket, err := l.c.GetBucket(l.Bucket())
	if err != nil {
		return err
	}
	slog.Debug("unlocking %s", slog.String("key", l.key))
	return bucket.Delete(l.key)
}
