package handlers

import (
	"time"

	"github.com/google/uuid"
)

func parseDate(s string) (time.Time, error) {
	return time.Parse(time.DateOnly, s)
}

func generateUUID() string {
	return uuid.New().String()
}
