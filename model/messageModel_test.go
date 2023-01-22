package model

import (
	"fmt"
	"testing"
)

func TestSelectMessagesByUserId(t *testing.T) {
	messages,err:=SelectMessagesByUserId(2,3)
	fmt.Printf("%v\n",messages)
	fmt.Printf("%v\n",err)

}