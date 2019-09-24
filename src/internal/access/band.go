package access

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/brucebales/bandmanager-backend/src/internal/structs"
)

//CreateBandJob is the struct passed into the band creation channel
type CreateBandJob struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	User        structs.User
}

//EditBandJob is the struct passed into the band edit channel
type EditBandJob struct {
	ID          int    `json:"band_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	User        structs.User
}

//CreateBandWorker creates new bands in the DB and adds the creator as the first member.
//This function runs in async using the CreateBandJob struct.
func CreateBandWorker(j <-chan CreateBandJob, wg *sync.WaitGroup, db *sql.DB) error {
	for job := range j {
		res, err := db.Exec("INSERT INTO prim.bands(name, description) VALUES(?, ?);", job.Name, job.Description)
		if err != nil {
			fmt.Println("Could not create band: ", err)
		}
		bandid, err := res.LastInsertId()
		if err != nil {
			fmt.Println("Could not find last inserted band ID")
		}
		_, err = db.Exec("INSERT INTO prim.bands_members(user_id, band_id, name, role, acl) VALUES(?, ?, ?, ?, ?);", job.User.ID, bandid, job.User.Name, "Founder", 4)
		if err != nil {
			fmt.Println("Could not create band: ", err)
		}
	}
	wg.Done()
	return nil
}

//EditBandWorker updates band information using info received in the EditBandJob channel
func EditBandWorker(j <-chan EditBandJob, wg *sync.WaitGroup, db *sql.DB) error {
	for job := range j {
		//Check that user requesting change is permitted to do so
		//TODO: Split these ACL checks into their own function
		members := make([]structs.Member, 0)
		rows, err := db.Query("SELECT user_id, name, role, acl FROM prim.bands_members WHERE band_id = ?", job.ID)
		for rows.Next() {
			memb := structs.Member{}
			err = rows.Scan(&memb.UserID, &memb.Name, &memb.Role, &memb.ACL)
			if err != nil {
				return err
			}
			members = append(members, memb)
		}
		var permitted = false
		for _, m := range members {
			if m.UserID == job.User.ID && m.ACL >= 3 {
				permitted = true
			}
		}
		if permitted {
			_, err = db.Exec("UPDATE prim.bands SET name=?, description=? WHERE id = ?;", job.Name, job.Description, job.ID)
			if err != nil {
				fmt.Println("Could not edit band: ", err)
			}
		} else {
			fmt.Println("User requested unauthorized update: ", job.User.ID, "Band ID: ", job.ID)
		}
	}
	wg.Done()
	return nil
}

//GetBandInfo fetches a band's info from it's ID
func GetBandInfo(bandID, userID int, db *sql.DB) (structs.Band, error) {

	band := structs.Band{}
	members := make([]structs.Member, 0)

	rows, err := db.Query("SELECT id, name, description FROM prim.bands WHERE id = ? LIMIT 1", bandID)
	if err != nil {
		return structs.Band{}, err
	}
	for rows.Next() {
		rows.Scan(&band.ID, &band.Name, &band.Description)
	}

	rows, err = db.Query("SELECT user_id, name, role, acl FROM prim.bands_members WHERE band_id = ?", bandID)
	for rows.Next() {
		memb := structs.Member{}
		err = rows.Scan(&memb.UserID, &memb.Name, &memb.Role, &memb.ACL)
		if err != nil {
			return structs.Band{}, err
		}
		members = append(members, memb)
	}
	band.Members = members

	var permitted = false
	for _, m := range members {
		if m.UserID == userID && m.ACL >= 1 {
			permitted = true
		}
	}
	if !permitted {
		return structs.Band{}, fmt.Errorf("User is not permitted to view this band's info")
	}

	return band, nil
}
