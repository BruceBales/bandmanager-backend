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

	/* --- Root Endpoint --- */
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Invalid path")
	})

	/* --- Auth Endpoints --- */
	http.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		sess, err := auth.Login(r.PostFormValue("password"), r.PostFormValue("email"), db)
		if err != nil {
			fmt.Println("Cannot login: ", err)
			w.WriteHeader(401)
			fmt.Fprintf(w, "Cannot login: %s", err)
		}
		session, err := json.Marshal(sess)
		if err != nil {
			//Error hash included in Println so that it can be easily searched in server logs
			fmt.Println("[b4bb70ea0a801c3d0286c2c4678b01a36a28a5a4e4e36d1a1b95a4b42fed2ffd] Cannot unmarshall JSON: ", err)
			w.WriteHeader(500)
			//Error hash returned in HTTP response so that whoever sees it can report it but won't accidentally see
			//an internal error message from Go.
			fmt.Fprintf(w, "Internal server error: b4bb70ea0a801c3d0286c2c4678b01a36a28a5a4e4e36d1a1b95a4b42fed2ffd")
		}
		fmt.Fprintf(w, string(session))
	})
	http.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		err := auth.CreateUSer(r.PostFormValue("name"), r.PostFormValue("email"), r.PostFormValue("password"), db)
		switch err {
		case nil:
			fmt.Fprintf(w, "Success!")
			break
		default:
			fmt.Println("[a9e64a7b779c290fe1918e819f59a560720cb757b433e20fc16042555acecd35] - Cannot create user: ", err)
			w.WriteHeader(400)
			fmt.Fprintf(w, "Cannot create user: a9e64a7b779c290fe1918e819f59a560720cb757b433e20fc16042555acecd35")
		}
	})

	err = http.ListenAndServe(":1226", nil)
	if err != nil {
		fmt.Println("HTTP Serve Error: ", err)
	}

}
