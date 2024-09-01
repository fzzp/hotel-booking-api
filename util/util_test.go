package util

import (
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	fmt.Println(MD5(MD5("87654321"))) // 1ba4413ca86ad65f676579cdf83d6752
}
