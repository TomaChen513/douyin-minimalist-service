package model

import (
	"fmt"
	"testing"
)

func TestGetTableRelation(t *testing.T) {
	flag:=GetTableRelation(2,3)
	fmt.Printf("%v\n", flag)
	flag=GetTableRelation(3,2)
	fmt.Printf("%v\n", flag)
	flag=GetTableRelation(1,2)
	fmt.Printf("%v\n", flag)
}