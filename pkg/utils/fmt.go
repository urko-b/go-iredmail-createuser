package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

// LogJSONObject logs an object as a JSON object.
func LogJSONObject(msg string, obj interface{}) {
	str, err := getJsonString(obj)
	if err != nil {
		fmt.Printf("l.getJsonString: %s\n", err)
		return
	}

	fmt.Printf("%s json_object: %s\n", msg, str)
}

func getJsonString(e interface{}) (string, error) {
	if e == nil {
		return "", errors.New("entity cannot be nil")
	}

	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
