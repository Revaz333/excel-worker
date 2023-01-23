package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Md5(data string) string {
	response := md5.Sum([]byte(data))
	return fmt.Sprintf("%s", hex.EncodeToString(response[:]))
}

func Call(fn interface{}, args []interface{}) interface{} {
	method := reflect.ValueOf(fn)
	var inputs []reflect.Value

	for _, v := range args {
		inputs = append(inputs, reflect.ValueOf(v))
	}

	return method.Call(inputs)[0]
}

func GenerateToken(length int, runeText string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return Md5(string(b) + runeText)
}

func Replace(data string, to string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	safe := reg.ReplaceAllString(data, to)
	safe = strings.ToLower(strings.Trim(safe, to))

	return safe
}
