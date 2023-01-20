package model

import (
	"fmt"
	"testing"
)


func TestInsertFavourite(t *testing.T) {
	success:=InsertFavourite(TableFavourite{UserId: 2,VideoId: 3})
	fmt.Printf("%+v\n", success)
	success=InsertFavourite(TableFavourite{UserId: 2,VideoId: 3})
	fmt.Printf("%+v\n", success)
}

func TestDeleteFavourite(t *testing.T) {
	// 无记录也可以删除
	success:=DeleteFavourite(TableFavourite{UserId: 2,VideoId: 3})
	fmt.Printf("%+v\n", success)
	success=DeleteFavourite(TableFavourite{UserId: 2,VideoId: 3})
	fmt.Printf("%+v\n", success)
}