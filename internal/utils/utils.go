package utils

import (
	"log"
	"os"
	"strconv"
)

// GetStringFromEnv to read string value from environment
func GetStringFromEnv(key, fallback string) string {
	if val, found := os.LookupEnv(key); found {
		log.Printf("reading key=%s from os environment found\n", key)
		return val
	}
	log.Printf("reading key=%s from os environment not found. returning fallback value\n", key)
	return fallback
}

// GetNumberFromEnv to read number value from environment
func GetNumberFromEnv(key string, fallback int) int {
	if val, found := os.LookupEnv(key); found {
		log.Printf("reading key=%s from os environment found", key)
		num, err := strconv.Atoi(val)
		if err != nil {
			log.Printf("something went wrong in reading %s key\n", key)
			return fallback
		}
		return num
	}
	log.Printf("reading key=%s from os environment not found. returning fallback value\n", key)
	return fallback
}
