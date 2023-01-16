package configs

import (
	"fmt"
	"testing"
)

func TestConf(t *testing.T) {
	conf, err := read()
	if err != nil {
		t.Log(err)
	}
	fmt.Println(conf)
}
