package service

import (
	"fmt"
	"testing"
)

func TestReleaseToken(t *testing.T) {
	token,_:=ReleaseToken(User{Id: 1})
	fmt.Println(token)
}

func TestParseToken(t *testing.T) {
	token,_:=ReleaseToken(User{Id: 123})
	claims,_:=ParseToken(token)
	fmt.Println(claims.UserId)
}