package tools

import "fmt"

func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}

func requireString(args map[string]any, key string) (string, error) {
	val, ok := args[key]
	if !ok {
		return "", fmt.Errorf("missing %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%s must be string", key)
	}
	return str, nil
}
