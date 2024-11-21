package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
)

const (
	simCliPrefix = "sim-cli-managed"
)

// CreateImage will build a new image using the predefined support-bundle-kit baseImage and layer it with the actual
// support bundle in /bundle directory. This can subsequently be loaded into the simulator
func (c *Client) CreateImage(instanceName string, bundlePath string, baseImage string) error {

	imageName := fmt.Sprintf("%s:%s", simCliPrefix, instanceName)
	contextTar, err := BuildContextTar(bundlePath, baseImage)
	if err != nil {
		return err
	}

	imageBuildResponse, err := c.APIClient.ImageBuild(c.ctx, bytes.NewReader(contextTar.Bytes()), types.ImageBuildOptions{
		Tags: []string{imageName},
		Labels: map[string]string{
			bundleNameKey: instanceName,
		},
	})

	if err != nil {
		return err
	}

	return readResponse(imageBuildResponse.Body)
}

// FindImage attempts to find image for a given instanceName by filtering on labels added
// to image during the image generation process
func (c *Client) FindImages(instanceName string) ([]image.Summary, error) {
	imageName := fmt.Sprintf("%s:%s", simCliPrefix, instanceName)
	filters := filters.NewArgs(filters.KeyValuePair{Key: "reference", Value: imageName})
	return c.APIClient.ImageList(c.ctx, image.ListOptions{
		Filters: filters,
	})
}

// RemoveImages removes images associated with instanceName
func (c *Client) RemoveImages(instanceName string) error {
	images, err := c.FindImages(instanceName)
	if err != nil {
		return nil
	}

	for _, v := range images {
		resp, err := c.APIClient.ImageRemove(c.ctx, v.ID, image.RemoveOptions{})
		if err != nil {
			return fmt.Errorf("error removing image %s: %v", v.ID, err)
		}

		for _, v := range resp {
			logrus.Infof("removed image: %v", v)
		}
	}

	return nil
}

// readResponse attempts to tidy up response messages
func readResponse(resp io.ReadCloser) error {
	defer resp.Close()
	reader := bufio.NewReader(resp)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		msg := &jsonmessage.JSONMessage{}
		err = json.Unmarshal(line, msg)
		if err != nil {
			return fmt.Errorf("error unmarshalling json: %v", err)
		}

		if msg.Error != nil {
			logrus.Error(msg.Error)
		}

		if msg.Aux != nil {
			continue
		}

		if msg.Stream != "" && msg.Stream != "\n" {
			logrus.Info(msg.Stream)
		}
	}
	return nil
}
