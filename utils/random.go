package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alp = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alp)

	for i := 0; i < n; i++ {
		c := alp[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%v@gmail.com", RandomString(6))

}

func RandomPhone() string {
	return fmt.Sprintf("0%v", RandomInt(100000000, 999999999))
}

func RandomName() string {
	return fmt.Sprintf("%v %v", RandomString(6), RandomString(6))
}

func RandomAddress() string {
	return fmt.Sprintf("%v %v", RandomString(6), RandomString(6))
}


func RandomFileName() string {
	return fmt.Sprintf("%v.jpg", RandomString(6))
}
