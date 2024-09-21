package internal

import (
	"context"
	"fmt"
)

type Notifier interface {
	Send(ctx context.Context, message string) error
}

type StdoutNotifier struct{}

func (n StdoutNotifier) Send(ctx context.Context, message string) error {
	fmt.Println("---")
	fmt.Printf("\n: %s :", message)
	fmt.Println("---")
	fmt.Println("")
	return nil
}
