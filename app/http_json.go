package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var (
	instance *http.Client
	once     sync.Once
)

func getClient() *http.Client {
	once.Do(func() {
		instance = &http.Client{
			Timeout: 10 * time.Second,
		}
	})
	return instance
}

func Post[T any](url string, data any) (*T, error) {
	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	client := getClient()
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(serialized))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 204 || resp.ContentLength == 0 {
		return nil, err
	}
	defer resp.Body.Close()
	var res T
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func HttpGet[T any](url string) (*T, error) {
	resp, err := getClient().Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data T
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
