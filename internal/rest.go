package internal

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"dominikdev.com/dogger/config"
	"dominikdev.com/dogger/internal/api"
)

type RESTNotifier struct {
	Webhooks      []config.Hook
	BodyProviders map[string]api.BodyProvider
}

func (m *RESTNotifier) Send(ctx context.Context, message string) error {
	var enabledHooks []config.Hook

	for _, hook := range m.Webhooks {
		if hook.Enabled {
			enabledHooks = append(enabledHooks, hook)
		}
	}

	var wg sync.WaitGroup
	resultChan := make(chan api.HttpResult, len(enabledHooks))

	for _, hook := range enabledHooks {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			body, err := m.buildBody(hook, message)
			if err == nil {
				performPOST(url, body, resultChan)
			} else {
				fmt.Printf("couldn't send message: %v\n", err)
			}
		}(hook.Url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		HandleResult(result)
	}

	return nil
}

func (m *RESTNotifier) buildBody(hook config.Hook, message string) ([]byte, error) {
	provider, exists := m.BodyProviders[hook.Name]
	if !exists {
		return []byte{}, fmt.Errorf("unsupported webhook: %s", hook.Name)
	}

	return provider.Provide(message)
}

func HandleResult(result api.HttpResult) {
	if result.Error != nil {
		fmt.Printf("Error with request to %s: %v\n", result.Url, result.Error)
		return
	}

	switch result.StatusCode {
	case http.StatusOK:
		fmt.Printf("Success! Request to %s returned status 200 OK.\n", result.Url)
	case http.StatusInternalServerError:
		fmt.Printf("Server error! Request to %s returned status 500.\n", result.Url)
	default:
		fmt.Printf("Request to %s returned status %d.\n", result.Url, result.StatusCode)
	}
}

func performPOST(url string, payload []byte, result chan<- api.HttpResult) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		result <- api.HttpResult{Url: url, StatusCode: 0, Error: err}
		return
	}

	req.Header.Set("Content-Type", "application/json;utf-8")
	resp, err := client.Do(req)
	if err != nil {
		result <- api.HttpResult{Url: url, StatusCode: 0, Error: err}
		return
	}

	result <- api.HttpResult{Url: url, StatusCode: resp.StatusCode, Error: nil}
}
