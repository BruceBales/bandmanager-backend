package main

import (
	"fmt"
	"sync"

	"github.com/brucebales/bandmanager-backend/src/internal/access"
	"github.com/brucebales/bandmanager-backend/src/internal/api"
	"github.com/brucebales/bandmanager-backend/src/internal/dao"
)

func main() {

	//Initialize DB connection
	db, err := dao.NewDB()
	if err != nil {
		fmt.Println("Failed to connect to db: ", err)
	}
	defer db.Close()

	//Create async waitgroup
	wg := new(sync.WaitGroup)

	//Create band creation channel
	createBandChan := make(chan access.CreateBandJob)
	//Spawn band creation worker
	wg.Add(1)
	go access.CreateBandWorker(createBandChan, wg, db)

	api.Serve(createBandChan, db)
}
