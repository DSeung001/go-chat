// Package utils contains functions to be used across the application
package utils

import (
	"encoding/json"
	"log"
)

var logFn = log.Panic

func HandleErr(err error) {
	if err != nil {
		logFn(err)
	}
}

func StructToBytes(data interface{}) []byte {
	bytes, err := json.Marshal(data)
	HandleErr(err)
	return bytes
}
