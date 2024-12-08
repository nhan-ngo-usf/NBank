package util

import (
	"math/rand"
	"strings"
)

// RandomInt generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	randomNum := rand.Int63n(max - min + 1) + min
	return randomNum
}

// RandomString create a random string of length n
var alphabet string = "abcdefghjklmnopqrstuvwyz"
func RandomString(n int) string {
	bound := len(alphabet)
	var strBuilder strings.Builder
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(bound)]
		strBuilder.WriteByte(c)
	}
	return strBuilder.String()
}

func RandomUserName() string {
	var randUserName string = RandomString(6)
	return randUserName
}

func RandomBalance() int64 {
	var randBalance int64 = RandomInt(10,20000)
	return randBalance
}

func RandomCurrency() string {
	var CurrencyList = []string {"USD", "EUR", "POUND"}
	numCurrency := len(CurrencyList)
	var randCurrency string = CurrencyList[rand.Intn(numCurrency)]
	return randCurrency
}