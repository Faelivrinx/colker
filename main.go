package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"dominikdev.com/dogger/config"
	"dominikdev.com/dogger/internal"
	"dominikdev.com/dogger/internal/api"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	nativeDockerClient "github.com/docker/docker/client"
)

func New() (*nativeDockerClient.Client, error) {
	dockerClient, error := nativeDockerClient.NewClientWithOpts(nativeDockerClient.FromEnv, nativeDockerClient.WithAPIVersionNegotiation())
	return dockerClient, error
}

func main() {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("couldn't read config file: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()

	dockerClient, error := New()
	if error != nil {
		fmt.Printf("Error creating docker client %v", error)
		return
	}
	eventFilter := filters.NewArgs()
	eventFilter.Add("type", "container")
	eventFilter.Add("action", "create")
	eventFilter.Add("action", "destroy")

	var notifier internal.Notifier
	notifier = &internal.RESTNotifier{Webhooks: config.Webhooks, BodyProviders: map[string]api.BodyProvider{
		"ms-teams": &api.MsTeamsProvider{},
	}}
	listenerManager := internal.NewHealthCheckManager(config.Messages.FinalMessage, notifier)
	// go listenerManager.DisplayState()

	messages, _ := dockerClient.Events(ctx, events.ListOptions{
		Filters: eventFilter,
	})

	notifier.Send(ctx, "Docker dogger started for listening!")

	for message := range messages {
		var (
			shouldBeProcessed bool
			statusUrl         string
		)
		containerName := message.Actor.Attributes["name"]
		image := message.Actor.Attributes["image"]

		for _, container := range config.Containers {
			if strings.ToLower(container.Name) == strings.ToLower(containerName) {
				statusUrl = container.StatusURL
				shouldBeProcessed = true
			}
		}

		if !shouldBeProcessed {
			continue
		}

		if message.Action == "start" {
			var version string
			regex := regexp.MustCompile(`[^:]+:([\w.-]+)$`)
			match := regex.FindStringSubmatch(image)
			if len(match) > 0 {
				version = match[1]
			} else {
				version = "unknown"
			}

			notifier.Send(ctx, fmt.Sprintf(config.Messages.StartMessage, containerName, version))
			if statusUrl != "" {
				listenerManager.RegisterListener(containerName, statusUrl)
			}
		}

		if message.Action == "die" {
			notifier.Send(ctx, fmt.Sprintf(config.Messages.StopMessage, containerName))
			listenerManager.UnregisterListener(containerName)
		}
	}
}
