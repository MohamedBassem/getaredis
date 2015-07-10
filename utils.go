package getaredis

import (
	"fmt"
	"math/rand"
)

func generateDockerAddress(ip, username, password string) string {
	return fmt.Sprintf("tcp://%v:%v@%v:2377", username, password, ip)
}

func generateRandomString(length int) string {
	var ret string
	runes := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	for i := 0; i < length; i++ {
		ret += string(runes[rand.Intn(len(runes))])
	}
	return ret
}
