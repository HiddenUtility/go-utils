package jsonutils

import (
	"encoding/json"
	"fmt"
	"os"
)

// ToJSON converts any Go value to a JSON string.
func ToJSON[T any](value T) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", JsonError{Message: fmt.Sprintf("failed to marshal JSON: %v", err)}
	}
	return string(bytes), nil
}

// FromJSON converts a JSON string to any Go value.
func FromJSON[T any](data string) (T, error) {
	var value T
	err := json.Unmarshal([]byte(data), &value)
	if err != nil {
		return value, JsonError{Message: fmt.Sprintf("failed to unmarshal JSON: %v", err)}
	}
	return value, nil
}

// ToJSONPretty converts any Go value to a formatted JSON string.
func ToJSONPretty[T any](value T) (string, error) {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", JsonError{Message: fmt.Sprintf("failed to marshal JSON pretty: %v", err)}
	}
	return string(bytes), nil
}

// WriteFile writes any Go value to a JSON file.
func WriteFile[T any](path string, value T) error {
	data, err := ToJSON(value)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0644)
}

// WriteFilePretty writes any Go value to a formatted JSON file.
func WriteFilePretty[T any](path string, value T) error {
	data, err := ToJSONPretty(value)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0644)
}

// ReadFile reads a JSON file and unmarshals it into the target.
func ReadFile[T any](path string) (T, error) {
	var value T
	bytes, err := os.ReadFile(path)
	if err != nil {
		return value, JsonError{Message: fmt.Sprintf("failed to read file: %v", err)}
	}
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return value, JsonError{Message: fmt.Sprintf("failed to unmarshal JSON from file: %v", err)}
	}
	return value, nil
}
