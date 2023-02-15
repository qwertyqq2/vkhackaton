package values

import (
	"fmt"
	"testing"
)

func TestLoop(t *testing.T) {
	var i = 0

	for ; i != 10; i++ {
		fmt.Println(i)
	}
}
