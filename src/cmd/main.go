package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brucebales/bandmanager-backend/src/internal/auth"
	"github.com/brucebales/bandmanager-backend/src/internal/config"
	"github.com/brucebales/bandmanager-backend/src/internal/dao"
)

func main() {

	conf := config.GetConfig()

	Database := dao.Database{
		Driver: "mysql",
		DS:     fmt.Sprintf("%s:%s@tcp(%s:%s)/prim", conf.MysqlUser, conf.MysqlPass, conf.MysqlHost, conf.MysqlPort),
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
			w.WriteHeader(401)
			fmt.Fprintf(w, "Cannot login: %s", err)
		}
		session, err := json.Marshal(sess)
		if err != nil {
			fmt.Println("Cannot unmarshall JSON: ", err)
			w.WriteHeader(500)
			//Using Sha256 hash of error description to mark where error is happening in code.
			//This allows me to quickly itentify where an error is happening without
			//Giving too many details to anyone else who might be using the API.
			fmt.Fprintf(w, "Internal server error: b4bb70ea0a801c3d0286c2c4678b01a36a28a5a4e4e36d1a1b95a4b42fed2ffd")
		}
		fmt.Fprintf(w, string(session))
	})

	err = http.ListenAndServe(":1226", nil)
	if err != nil {
		fmt.Println("HTTP Serve Error: ", err)
	}

}
