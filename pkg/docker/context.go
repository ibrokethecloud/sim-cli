package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

func GetCurrentContext(ctx context.Context) error {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error fetching client: %v", err)
	}
	sysInfo, err := apiClient.Info(ctx)
	if err != nil {
		return fmt.Errorf("error fetching system info: %v", err)
	}
	fmt.Println(sysInfo)
	return nil
}
