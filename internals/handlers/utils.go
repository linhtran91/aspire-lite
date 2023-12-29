package handlers

import (
	"aspire-lite/internals/constants"
	"net/url"
	"strconv"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var defaultAlphabel = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func parseDate(s string) (time.Time, error) {
	return time.Parse(time.DateOnly, s)
}

func generateUUID() string {
	id, _ := gonanoid.Generate(defaultAlphabel, constants.LengthOfID)
	return id
}

func getLimitOffset(values url.Values) (int, int) {
	page := getValueFromUrl(values, "page", constants.DefaultPage)
	size := getValueFromUrl(values, "size", constants.DefaultSize)
	if size >= constants.MaximumSize {
		size = constants.DefaultSize
	}
	if page == 1 {
		return size, 0
	}

	return size, (page-1)*size + 1
}

func getValueFromUrl(values url.Values, key string, defaultVal int) int {
	val := values.Get(key)
	if val == "" {
		return defaultVal
	}
	t, err := strconv.Atoi(val)
	if err != nil || t == 0 {
		return constants.DefaultPage
	}
	return t

}
