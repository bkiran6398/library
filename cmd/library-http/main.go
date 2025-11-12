package main

import (
	"fmt"

	"github.com/bkiran6398/library/internal/config"
)

func main() {
	configuration, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	fmt.Println(configuration)
}
