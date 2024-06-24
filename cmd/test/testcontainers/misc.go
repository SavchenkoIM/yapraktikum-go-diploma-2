package testcontainers

import (
	"fmt"
	"github.com/testcontainers/testcontainers-go"
)

type CustomLogConsumer struct{}

func (c CustomLogConsumer) Accept(log testcontainers.Log) {
	fmt.Println(string(log.Content))
}
