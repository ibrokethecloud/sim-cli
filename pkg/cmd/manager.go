package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ibrokethecloud/sim-cli/pkg/kubeconfig"
	"github.com/sirupsen/logrus"
)

const (
	defaultSimKubeConfigPath = ".sim/admin.kubeconfig"
	defaultKubeConfigPath    = "/root/.sim/admin.kubeconfig"
)

// PreFlightChecks ensures that
func (s *Simulator) PreFlightChecks() error {
	// check bundlePath exists
	bundleInfo, err := os.Stat(s.BundlePath)
	if err != nil {
		return fmt.Errorf("error checking bundle path %s: %w", s.BundlePath, err)
	}

	// ensure bundlePath is not a directory
	if bundleInfo.IsDir() {
		return fmt.Errorf("bundlePath needs to be location of zip file, current path %s is a directory", s.BundlePath)
	}

	// check if a container is already running
	containers, err := s.DockerClient.FindRunningContainer(s.Name)
	if err != nil {
		return fmt.Errorf("error listing running containers: %w", err)
	}

	if len(containers) != 0 {
		var ids []string
		for _, v := range containers {
			ids = append(ids, v.ID)
		}
		return fmt.Errorf("found containers with ID's %v already running, please stop existing containers or use a different name argument", ids)
	}

	return nil
}

// CreateNewInstall will deploy a new instance of the simulator using the support bundle
func (s *Simulator) CreateNewInstance() error {
	if err := s.DockerClient.CreateImage(s.Name, s.BundlePath, s.Image); err != nil {
		return fmt.Errorf("error creating new sim image: %w", err)
	}

	//run newly create image
	if err := s.DockerClient.RunContainer(s.Name, s.BundlePath); err != nil {
		return fmt.Errorf("error running new image: %w", err)
	}

	containers, err := s.DockerClient.FindRunningContainer(s.Name)
	if err != nil {
		return fmt.Errorf("error listing running containers: %w", err)
	}

	if len(containers) != 1 {
		return fmt.Errorf("expected to find only 1 running container but found %d", len(containers))
	}

	s.Port = int(containers[0].Ports[0].PublicPort)
	logrus.WithField("name", s.Name).Infof("simulator instance exposed on port %d", s.Port)
	return nil
}

// ListInstances will report the details of currently running sim instances
func (s *Simulator) ListInstances() error {
	return s.DockerClient.FindAllSimManagedInstances()
}

func (s *Simulator) ExportKubeConfig() error {
	logrus.Infof("exporting kubeconfig for instance %s", s.Name)
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error fetching home directory: %w", err)
	}

	kubeConfigPath := filepath.Join(home, defaultSimKubeConfigPath)

	contents, err := s.DockerClient.ReadFile(s.Name, defaultKubeConfigPath)
	if err != nil {
		return fmt.Errorf("error fetching kubeconfig from container %s: %w", s.Name, err)
	}

	endpoint, port, err := s.DockerClient.QueryExposedMapping(s.Name)
	if err != nil {
		return err
	}
	err = kubeconfig.AddContext(kubeConfigPath, s.Name, endpoint, port, contents)
	if err != nil {
		return fmt.Errorf("error adding context for %s to kubeconfig: %w", s.Name, err)
	}
	logrus.Infof("exported kubeconfig to context %s", s.Name)
	return nil
}

func (s *Simulator) RemoveInstance() error {
	logrus.Infof("removing instance %s", s.Name)
	if err := s.DockerClient.StopContainer(s.Name); err != nil {
		return fmt.Errorf("error removing container %s: %w", s.Name, err)
	}

	logrus.Infof("removing image for instance %s", s.Name)
	if err := s.DockerClient.RemoveImages(s.Name); err != nil {
		return fmt.Errorf("error removing image for instance %s: %w", s.Name, err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error fetching home directory: %w", err)
	}
	logrus.Infof("removing context for instance %s", s.Name)
	kubeConfigPath := filepath.Join(home, defaultSimKubeConfigPath)
	return kubeconfig.RemoveContext(kubeConfigPath, s.Name)
}
