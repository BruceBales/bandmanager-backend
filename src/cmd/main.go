package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brucebales/bandmanager-backend/src/internal/auth"
	"github.com/brucebales/bandmanager-backend/src/internal/dao"
)

func main() {

	conf := getConfig()

	Database := dao.Database{
		Driver: "mysql",
		DS:     fmt.Sprintf("%s:%s@tcp(%s:%s)/prim", conf.mysqlUser, conf.mysqlPass, conf.mysqlHost, conf.mysqlPort),
	}

	db, err := Database.New()
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Invalid path")
	})
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		sess, err := auth.Login(r.PostFormValue("password"), r.PostFormValue("email"), db)
		if err != nil {
			fmt.Println("Cannot login: ", err)
		}
		session, err := json.Marshal(sess)
		if err != nil {
			fmt.Println("Cannot unmarshall JSON: ", err)
		}
		fmt.Fprintf(w, string(session))
	})

	err = http.ListenAndServe(":1226", nil)
	if err != nil {
		fmt.Println("HTTP Serve Error: ", err)
	}

}
