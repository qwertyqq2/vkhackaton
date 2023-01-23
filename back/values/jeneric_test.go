package values

import (
	"fmt"
	"testing"
)

func PrintSlice[T any](s []T) {
	for _, val := range s {
		fmt.Println(val)
	}
}

type Mytype struct {
	a uint

	b uint
}

func Test(t *testing.T) {
}
