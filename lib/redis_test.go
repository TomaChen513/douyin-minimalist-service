package lib

import (
	"fmt"
	"testing"
)

func TestSetKey(t *testing.T) {
	SetKey("1223","1223",10)
	SetKey("1456","145",10)

}

func TestGetKey(t *testing.T) {
	fmt.Println(GetKey("1223"))
	fmt.Println(GetKey("1456"))
}