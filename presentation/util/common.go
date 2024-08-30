// The presentation layer contains all resources concerned with creating an application interface
// This file contains common code that apply to any type of interface (CLI, native GUI or web GUI or whatever..)
package util

import (

	"math/rand"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


func GenerateRandomNumber(min int , max int) int {
	return rand.Intn(max-min) + min
}
