package kv

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"paas/util"
	"time"
)

type Lock struct {
	key     string
	timeout time.Duration
	c       *Client
}

func (c *Client) NewLock(key string, timeout time.Duration) *Lock {
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
				fmt.Printf("waiting for lock %s\n", l.key)
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
	fmt.Printf("unlocking %s\n", l.key)
	return bucket.Delete(l.key)
}
