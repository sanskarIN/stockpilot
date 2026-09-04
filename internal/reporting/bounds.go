package reporting

import (
	"fmt"
	"strconv"
	"strings"
)

type Bounds struct {
	Limit  int
	Offset int
}

func NewBounds(limit, offset, defaultLimit, maxLimit int) (Bounds, error) {
	if defaultLimit <= 0 || maxLimit < defaultLimit {
		return Bounds{}, fmt.Errorf("invalid bounds configuration")
	}
	if limit == 0 {
		limit = defaultLimit
	}
	if limit < 1 || limit > maxLimit {
		return Bounds{}, fmt.Errorf("limit must be between 1 and %d", maxLimit)
	}
	if offset < 0 {
		return Bounds{}, fmt.Errorf("offset must not be negative")
	}
	return Bounds{Limit: limit, Offset: offset}, nil
}

func ParsePositiveInt(value string, fallback int) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("value must be a non-negative integer")
	}
	return parsed, nil
}
