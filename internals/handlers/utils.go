package handlers

import (
	"aspire-lite/internals/constants"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func getPageAndSize(values url.Values) (int, int) {
	page := getValueFromUrl(values, "page", constants.DefaultPage)
	size := getValueFromUrl(values, "size", constants.DefaultSize)
	if size >= constants.MaximumSize {
		size = constants.DefaultSize
	}

	return size, (page - 1) * size
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

func generateToken(hmacSampleSecret string, customerID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": customerID,
		"nbf":         time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	mySigningKey := []byte(hmacSampleSecret)
	return token.SignedString(mySigningKey)
}

func parseToken(tokenString, hmacSampleSecret string) (string, error) {
	mySigningKey := []byte(hmacSampleSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return fmt.Sprintf("%v", claims["customer_id"]), nil
	}
	return "", nil
}
