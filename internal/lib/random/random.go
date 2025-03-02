package random

import (
	"math/rand"
	"time"
)

// NewRandomString generates random string with given size.
func NewRandomString(size int) string {
	// создание псевдослучайного источника чисел
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	alf := []rune("0123456789" +
		"QWERTYUIOPASDFGHJKLZXCVBNM" +
		"qwertyuiopasdfghjklzxcvbnm")

	alias := make([]rune, size)

	for i := range alias {
		alias[i] = alf[rnd.Intn(len(alf))]
	}

	return string(alias)
}
