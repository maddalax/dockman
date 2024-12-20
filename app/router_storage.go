package app

import (
	"dockman/app/util/json2"
	"encoding/json"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
)

func ApplyBlocks(locator *service.Locator, blocks []RouteBlock) error {
	err := ValidateBlocks(blocks)

	if err != nil {
		return err
	}

	// Save block to database
	client := KvFromLocator(locator)
	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "route-table",
	})
	if err != nil {
		return err
	}
	_, err = bucket.Put("blocks", json2.SerializeOrEmpty(blocks))
	if err != nil {
		return err
	}

	ReloadConfig(locator)

	return nil
}

func ValidateBlocks(blocks []RouteBlock) error {
	for _, block := range blocks {
		err := ValidateBlock(&block)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRouteTable(locator *service.Locator) ([]RouteBlock, error) {
	client := KvFromLocator(locator)
	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "route-table",
	})
	if err != nil {
		return nil, err
	}

	data, err := bucket.Get("blocks")

	// empty table
	if err != nil {
		return make([]RouteBlock, 0), nil
	}
	var blocks []RouteBlock

	err = json.Unmarshal(data.Value(), &blocks)

	if err != nil {
		return nil, err
	}

	return blocks, nil
}
