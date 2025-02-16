package utils

import (
	//	"fmt"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generates a random string takes in aguement an int that is the length of the ouput string
func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

/*func main() {
	randomString := generateRandomString(8) // Here, we specify 8 to generate an 8-character string
	fmt.Println("Random 8-character string:", randomString)
}*/
