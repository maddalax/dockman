package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/maddalax/htmgo/framework/service"
	"paas/kv"
)

func ApplyBlocks(locator *service.Locator, blocks []RouteBlock) error {
	err := ValidateBlocks(blocks)

	if err != nil {
		return err
	}

	// Save block to database
	client := kv.GetClientFromLocator(locator)
	bucket, err := client.GetBucket("route-table")
	if err != nil {
		return err
	}
	_, err = bucket.Put("blocks", kv.MustSerialize(blocks))
	if err != nil {
		return err
	}

	ReloadConfig(locator)

	return nil
}

func ValidateBlocks(blocks []RouteBlock) error {
	// Validate blocks

	for i, block := range blocks {
		if block.Hostname == "" {
			return errors.New(fmt.Sprintf("Hostname is required for block %d", i))
		}
		if block.ResourceId == "" {
			return errors.New(fmt.Sprintf("resource id is required for block %d", i))
		}
		// TODO ensure the resource exists
	}

	return nil
}

func GetRouteTable(locator *service.Locator) ([]RouteBlock, error) {
	client := kv.GetClientFromLocator(locator)
	bucket, err := client.GetBucket("route-table")
	if err != nil {
		return nil, err
	}
	data, err := bucket.Get("blocks")
	if err != nil {
		return nil, err
	}
	var blocks []RouteBlock

	err = json.Unmarshal(data.Value(), &blocks)

	if err != nil {
		return nil, err
	}

	return blocks, nil
}
