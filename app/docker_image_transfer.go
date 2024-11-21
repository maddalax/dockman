package app

import (
	"context"
	"dockside/app/logger"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

func (c *DockerClient) LoadImage(imageId string) error {
	if c.HasLatestImage(imageId) {
		logger.InfoWithFields("We have the latest image, skipping load", map[string]any{
			"imageId": imageId,
		})
		return nil
	}

	logger.InfoWithFields("We don't have the latest image, loading from store", map[string]any{
		"imageId": imageId,
	})

	store, err := KvFromLocator(c.locator).ImageStore()

	if err != nil {
		return errors.Wrap(err, "failed to get object store")
	}

	obj, err := store.Get(imageId)

	if err != nil {
		return errors.Wrap(err, "failed to get docker image")
	}

	defer obj.Close()
	_, err = c.cli.ImageLoad(context.Background(), obj, true)

	if err != nil {
		return errors.Wrap(err, "failed to load docker image from store")
	}

	return nil
}

func (c *DockerClient) HasLatestImage(imageId string) bool {
	imageInfo, _, err := c.cli.ImageInspectWithRaw(context.Background(), imageId)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false
		} else {
			logger.Error("Error inspecting image", err)
			return false
		}
	}

	// has the image, but is it the latest?
	store, err := KvFromLocator(c.locator).ImageStore()

	if err != nil {
		logger.Error("Failed to get object store", err)
		return false
	}

	buildId := store.GetBuildId(imageId)
	currentBuildId := imageInfo.Config.Labels["dockside.build.id"]

	logger.InfoWithFields("Checking docker image, if we have latest", map[string]interface{}{
		"imageId":        imageId,
		"newBuildId":     buildId,
		"currentBuildId": currentBuildId,
	})

	if buildId != currentBuildId {
		return false
	}

	return true
}

func (c *DockerClient) SaveImage(imageId string, buildId string) error {
	body, err := c.cli.ImageSave(context.Background(), []string{imageId})
	if err != nil {
		return err
	}
	defer body.Close()
	store, err := KvFromLocator(c.locator).ImageStore()
	if err != nil {
		return err
	}
	meta := &nats.ObjectMeta{
		Name: fmt.Sprintf("%s", imageId),
		Metadata: map[string]string{
			"buildId": buildId,
		},
	}
	_, err = store.Put(meta, body)
	if err != nil {
		return err
	}
	return nil
}
