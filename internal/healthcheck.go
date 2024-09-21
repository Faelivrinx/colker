package internal

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	interval = 2 * time.Second
	timeout  = 30 * time.Second
)

type HealthCheckManager struct {
	listeners map[string]context.CancelFunc
	notifier  Notifier
	mu        sync.Mutex
	message   string
}

func NewHealthCheckManager(message string, notifier Notifier) *HealthCheckManager {
	return &HealthCheckManager{
		notifier:  notifier,
		message:   message,
		listeners: make(map[string]context.CancelFunc),
	}
}

func (hcm *HealthCheckManager) DisplayState() {
	timer := time.NewTicker(5 * time.Second)

	defer timer.Stop()

	for range timer.C {
		fmt.Printf("\n---\nState: %v\n---\n", hcm.listeners)
	}
}

func (hcm *HealthCheckManager) UnregisterListener(container string) {
	hcm.mu.Lock()
	defer hcm.mu.Unlock()

	if cancel, exists := hcm.listeners[container]; exists {
		cancel()
		fmt.Printf("Unregistered listener for container: %v\n", container)
	} else {
		fmt.Printf("Listener for container %v doesn't not exist\n", container)
	}
}

func (hcm *HealthCheckManager) RegisterListener(container, statusURL string) {
	hcm.mu.Lock()
	defer hcm.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())

	if oldCancel, exists := hcm.listeners[container]; exists {
		oldCancel()
		fmt.Printf("Canceled old listener for container %v\n", container)
	}

	hcm.listeners[container] = cancel

	go hcm.startHealthCheck(ctx, container, statusURL)
}

func (hcm *HealthCheckManager) startHealthCheck(ctx context.Context, container, statusURL string) {
	ticker := time.NewTicker(interval)

	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-context.Done():
			fmt.Printf("Container %s is in unhealthy condition", container)
			return
		case <-ticker.C:
			if hcm.isContainerHealthy(statusURL) {
				hcm.notifier.Send(ctx, fmt.Sprintf(hcm.message, container))
				return
			}
		}
	}
}

func (hcm *HealthCheckManager) isContainerHealthy(statusURL string) bool {
	resp, err := http.Get(statusURL)
	if err != nil {
		fmt.Printf("Error checking health: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
