package model

import (
	"fmt"
	"testing"

)



func TestGetVideosByUserId(t *testing.T) {
	user,err:=GetVideosByUserId(2)
	fmt.Printf("%+v  ", user)
	fmt.Printf("%v\n", err)
	user,err=GetVideosByUserId(3)
	fmt.Printf("%+v  ", user)
	fmt.Printf("%v\n", err)
	user,err=GetVideosByUserId(-1)
	fmt.Printf("%+v  ", user)
	fmt.Printf("%v\n", err)
}