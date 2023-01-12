package internal

import (
	"math"
	"math/rand"
)

const randomStep = 1

// RandomInt32MinIMaxN From min (included) to max (non-included)
func RandomInt32MinIMaxN(min int32, max int32) int32 {
	var d = max - min

	r := rand.Float64()*float64(d) + float64(min)

	return int32(r)
}

// RandomInt32MinMaxI From min (included) to max (included)
func RandomInt32MinMaxI(min int32, max int32) int32 {
	var d = max - min + randomStep

	r := rand.Float64()*float64(d) + float64(min)

	return int32(math.Floor(r))
}

// RandomInt32MinNMaxI From min (non-included) to max (included)
func RandomInt32MinNMaxI(min int32, max int32) int32 {
	return RandomInt32MinIMaxN(min+randomStep, max+randomStep)
}

// RandomInt32MinMaxN From min (non-included) to max (non-included)
func RandomInt32MinMaxN(min int32, max int32) int32 {
	return RandomInt32MinMaxI(min+randomStep, max-randomStep)
}

// ItoB Convert number != 0 to true, other false
func ItoB(number int32) bool {
	return number != 0
}
