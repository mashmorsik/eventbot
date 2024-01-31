package telegram

import (
	"strings"
	"time"
)

func isValidDate(dateStr string) bool {
	layout := "2006-01-02"

	_, err := time.Parse(layout, dateStr)

	return err == nil
}

func isValidTime(timeStr string) bool {
	layout := "15:04"

	_, err := time.Parse(layout, timeStr)

	return err == nil
}

func isValidFrequency(freq string) bool {
	options := []string{"once", "daily", "weekly", "monthly", "yearly"}

	frequency := strings.ToLower(freq)

	for _, o := range options {
		if frequency == o {
			return true
		}
	}

	return false
}
