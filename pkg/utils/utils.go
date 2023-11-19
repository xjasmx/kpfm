package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParsePortMapping(mapping string) (int, int, error) {
	parts := strings.Split(mapping, ":")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid port mapping format")
	}
	localPort, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid local port: %w", err)
	}
	podPort, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid pod port: %w", err)
	}
	return localPort, podPort, nil
}
