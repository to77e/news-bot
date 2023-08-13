package botkit

import (
	"encoding/json"
	"fmt"
)

func ParseJSON[T any](src string) (T, error) {
	var args T

	if err := json.Unmarshal([]byte(src), &args); err != nil {
		return *(new(T)), fmt.Errorf("json unmarshal: %w", err)
	}

	return args, nil
}
