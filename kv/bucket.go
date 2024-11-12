package kv

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/nats.go"
)

func Deserialize(data any, result any) error {
	config := &mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           &result,
	}

	config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		mapstructure.WeaklyTypedHook,
		mapstructure.TextUnmarshallerHookFunc(),
	)

	decoder, err := mapstructure.NewDecoder(config)

	if err != nil {
		return err
	}

	err = decoder.Decode(data)

	if err != nil {
		return err
	}

	return nil
}

func MustMapStringToStructure[T any](data string) *T {
	n := new(T)
	m := make(map[string]any)
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil
	}
	err = Deserialize(m, &n)
	if err != nil {
		return nil
	}
	return n
}

func AtomicPutMany(bucket nats.KeyValue, cb func(m map[string]any) error) error {
	m := make(map[string]any)
	err := cb(m)
	hasError := false

	if err != nil {
		return err
	}

	defer func() {
		if hasError {
			for k := range m {
				_ = bucket.Delete(k)
			}
		}
	}()

	for s, a := range m {
		_, err = bucket.Put(s, MustSerialize(a))
		if err != nil {
			hasError = true
			return err
		}
	}

	return err
}

func MustMapIntoMany[T any](client *Client, buckets <-chan string, keys ...string) ([]T, error) {
	n := make([]T, 0)
	for bucket := range buckets {
		b, err := client.GetBucket(bucket)
		if err != nil {
			continue
		}
		m, err := MustMapInto[T](b, keys...)
		if err != nil {
			return nil, err
		}
		n = append(n, *m)
	}
	return n, nil
}

func MustMapAllInto[T any](bucket nats.KeyValue) (*T, error) {
	n := new(T)
	m := make(map[string]string)
	k, err := bucket.ListKeys()

	if err != nil {
		return nil, err
	}
	for key := range k.Keys() {
		value, err := bucket.Get(key)
		if err != nil {
			return nil, err
		}
		v := string(value.Value())
		m[key] = v
	}

	err = Deserialize(m, &n)

	if err != nil {
		return nil, err
	}

	return n, nil
}

// MustMapInto fills in the struct with the values from the bucket, only if all the keys passed in exist in the bucket
func MustMapInto[T any](bucket nats.KeyValue, keys ...string) (*T, error) {
	n := new(T)
	m := make(map[string]string)
	for _, key := range keys {
		value, err := bucket.Get(key)
		if err != nil {
			return nil, err
		}
		v := value.Value()
		m[key] = string(v)
	}

	err := Deserialize(m, &n)

	if err != nil {
		return nil, err
	}

	return n, nil
}
