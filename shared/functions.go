package shared

import (
	"encoding/json"
	"fmt"
	"log"
)

func WriteError(format string, args ...interface{}) []byte {
	response := map[string]string{
		"error": fmt.Sprintf(format, args...),
	}

	if data, err := json.Marshal(response); err == nil {
		return data
	} else {
		log.Printf("Err: %s", err)
	}

	return nil
}

func WriteInfo(format string, args ...interface{}) []byte {
	response := map[string]string{
		"info": fmt.Sprintf(format, args...),
	}

	if data, err := json.Marshal(response); err == nil {
		return data
	} else {
		log.Printf("Err: %s", err)
	}

	return nil
}
