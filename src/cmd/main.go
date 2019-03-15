package main

import (
	"fmt"

	"github.com/bbales/bandmanager-backend/src/internal/dao"
)

func main() {
	db := dao.Database{
		Driver: "mysql",
		DS:     "",
	}
	conn, err := db.New()
	if err != nil {
		fmt.Println("!!Error: ", err)
	}
}
