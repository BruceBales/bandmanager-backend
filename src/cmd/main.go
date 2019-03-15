package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bbales/bandmanager-backend/src/internal/auth"
	"github.com/bbales/bandmanager-backend/src/internal/dao"
)

func main() {

	Database := dao.Database{
		Driver: "mysql",
		DS:     "test",
	}

	db, err := Database.New()

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
