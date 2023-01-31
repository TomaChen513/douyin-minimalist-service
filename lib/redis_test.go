package lib

import (
	"fmt"
	"testing"
)

func TestSetKey(t *testing.T) {
	SetKey("123","test",10)
}

func TestGetKey(t *testing.T) {
	fmt.Println(GetKey("123"))
}