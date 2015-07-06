package getaredis

import "math/rand"

func generateRandomString(length int) string {
	var ret string
	runes := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	for i := 0; i < length; i++ {
		ret += string(runes[rand.Intn(len(runes))])
	}
	return ret
}
