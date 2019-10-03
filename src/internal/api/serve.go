package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/brucebales/bandmanager-backend/src/internal/access"
	"github.com/brucebales/bandmanager-backend/src/internal/auth"
)

//Serve handles actual HTTP connections to the API
func Serve(channels access.WorkerChannels, responseChan <-chan access.Response, db *sql.DB) {
	/* --- Root Endpoint --- */
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Invalid path")
	})

	/* --- Business Logic Endpoints --- */
	//Create Band Endpoint.
	http.HandleFunc("/band/create", func(w http.ResponseWriter, r *http.Request) {
		crband := access.CreateBandJob{}
		d := json.NewDecoder(r.Body)
		d.Decode(&crband)

		sess := r.Header.Get("session_id")
		user, err := access.GetUser(sess, db)
		if err != nil {
			fmt.Fprintf(w, "Error getting user info")
		}
		crband.User = user
		channels.CreateBandChan <- crband
		fmt.Fprintf(w, "Success")
	})

	//Edit Band Endpoint
	http.HandleFunc("/band/edit", func(w http.ResponseWriter, r *http.Request) {

		editBand := access.EditBandJob{}
		d := json.NewDecoder(r.Body)
		d.Decode(&editBand)

		sess := r.Header.Get("session_id")
		user, err := access.GetUser(sess, db)
		if err != nil {
			fmt.Fprintf(w, "Error getting user info")
		}
		editBand.User = user

		channels.EditBandChan <- editBand
		fmt.Fprintf(w, "Success")
	})

	//Edit Member Endpoint
	http.HandleFunc("/band/member", func(w http.ResponseWriter, r *http.Request) {
		//Scan JSON into MemberJob struct
		membJob := access.MemberJob{}
		membJob.JobID = time.Now().Unix()
		d := json.NewDecoder(r.Body)
		err := d.Decode(&membJob)
		if err != nil {
			fmt.Fprintf(w, "Invalid JSON")
		}

		//Validate user session
		sess := r.Header.Get("session_id")
		user, err := access.GetUser(sess, db)
		if err != nil {
			fmt.Fprintf(w, "Error getting user info")
		}
		membJob.UserID = user.ID

		//Get user info for modified member
		membJob.Member, err = access.GetUserByID(membJob.MemberID, db)
		if err != nil {
			fmt.Fprintf(w, "Failed to get user ID: %v", membJob.MemberID)
			fmt.Println("Failed to get user ID ", membJob.MemberID, "Error: ", err)
		}
		channels.MemberChan <- membJob
		for resp := range responseChan {
			if resp.JobID == membJob.JobID {
				w.WriteHeader(resp.HTTPCode)
				fmt.Fprintf(w, resp.Message)
				return
			}
		}
	})

	//Get Band Endpoint
	http.HandleFunc("/band/info", func(w http.ResponseWriter, r *http.Request) {
		req := access.BandInfoRequest{}
		d := json.NewDecoder(r.Body)
		d.Decode(&req)

		sess := r.Header.Get("session_id")
		user, err := access.GetUser(sess, db)
		if err != nil {
			fmt.Fprintf(w, "Error getting user info")
		}

		bandInfo, err := access.GetBandInfo(req.BandID, user.ID, db)
		if err != nil {
			fmt.Println("Could not get band: ", err)
		}
		bandBytes, err := json.Marshal(bandInfo)
		fmt.Fprintf(w, string(bandBytes))
	})

	/* --- Auth Endpoints --- */
	//Login
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
	//Register
	http.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		err := auth.CreateUser(r.PostFormValue("name"), r.PostFormValue("email"), r.PostFormValue("password"), db)
		switch err {
		case nil:
			fmt.Fprintf(w, "Success!")
			break
		default:
			fmt.Println("[a9e64a7b779c290fe1918e819f59a560720cb757b433e20fc16042555acecd35] - Cannot create user: ", err)
			w.WriteHeader(400)
			if strings.Contains(fmt.Sprint(err), "Duplicate entry") {
				fmt.Fprintf(w, "Cannot create user: email already in use")
			} else {
				fmt.Fprintf(w, "Cannot create user: a9e64a7b779c290fe1918e819f59a560720cb757b433e20fc16042555acecd35")
			}
		}
	})

	err := http.ListenAndServe(":1226", nil)
	if err != nil {
		fmt.Println("HTTP Serve Error: ", err)
	}

}
