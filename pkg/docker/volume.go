package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/volume"
	"github.com/sirupsen/logrus"
)

func (c *Client) CreateVolume(name string) error {
	volume, err := c.APIClient.VolumeCreate(c.ctx, volume.CreateOptions{
		Name:   name,
		Driver: "local",
		Labels: map[string]string{
			bundleNameKey: "name",
		},
	})
	if err != nil {
		return fmt.Errorf("error during volume creation: %v", err)
	}
	logrus.WithField("volume", volume).Debug("volume created")
	return nil
}

func (c *Client) RemoveVolume(name string) error {
	return nil
}
