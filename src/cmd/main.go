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

	//Create needed job channels.
	//Currently async is only being used for functions
	//that modify database content
	createBandChan := make(chan access.CreateBandJob)
	editBandChan := make(chan access.EditBandJob)
	memJobChan := make(chan access.MemberJob)
	//Spawn band creation worker
	wg.Add(1)
	go access.CreateBandWorker(createBandChan, wg, db)
	wg.Add(1)
	go access.EditBandWorker(editBandChan, wg, db)
	wg.Add(1)
	go access.MemberWorker(memJobChan, wg, db)

	channels := access.WorkerChannels{
		CreateBandChan: createBandChan,
		EditBandChan:   editBandChan,
		MemberChan:     memJobChan,
	}

	//Start HTTP server
	api.Serve(channels, db)
}
