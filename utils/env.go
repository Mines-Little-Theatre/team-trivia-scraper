package utils

import (
	"log"
	"os"
)

// ReadEnv reads an environment variable, fatally logging an error message if it is not set
func ReadEnv(key string) string {
	var result string
	var ok bool
	if result, ok = os.LookupEnv(key); !ok {
		log.Fatalf("please set the %s environment variable", key)
	}
	return result
}
