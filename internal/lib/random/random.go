package random

import (
	"math/rand"
	"time"
)

func NewRandString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	cons := []rune{'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'z'}
	vowels := []rune{'a', 'e', 'i', 'o', 'u', 'y'}

	randAlias := make([]rune, size)

	for i := 0; i < size; i++ {
		randAlias[i] = coinFlip()

		switch randAlias[i] {
		case 'c':
			randAlias[i] = cons[rnd.Intn(len(cons))]
		case 'v':
			randAlias[i] = vowels[rnd.Intn(len(vowels))]
		}
	}
	return string(randAlias)
}

func coinFlip() rune {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	coin := []rune{
		'c',
		'v',
	}

	side := coin[rnd.Intn(len(coin))]

	return side
}
