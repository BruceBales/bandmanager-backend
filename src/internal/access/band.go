package access

import (
	"fmt"
	"sync"

	"github.com/brucebales/bandmanager-backend/src/internal/dao"
	"github.com/brucebales/bandmanager-backend/src/internal/structs"
)

type CreateBandJob struct {
	Name        string
	Description string
	User        structs.User
}

//CreateBandWorker creates new bands in the DB and adds the creator as the first member.
//This function runs in async using the CreateBandJob struct.
func CreateBandWorker(j <-chan CreateBandJob, wg *sync.WaitGroup) error {
	mysql, err := dao.NewMysql()
	if err != nil {
		return err
	}
	defer mysql.Close()

	for job := range j {
		res, err := mysql.Exec("INSERT INTO prim.bands(name, description) VALUES(?, ?);", job.Name, job.Description)
		if err != nil {
			fmt.Println("Could not create band: ", err)
		}
		bandid, err := res.LastInsertId()
		if err != nil {
			fmt.Println("Could not find last inserted band ID")
		}
		_, err = mysql.Exec("INSERT INTO prim.bands_members(user_id, band_id, name, role, acl) VALUES(?, ?, ?, ?, ?);", job.User.ID, bandid, job.User.Name, "Founder", 4)
		if err != nil {
			fmt.Println("Could not create band: ", err)
		}
	}
	wg.Done()
	return nil
}

//GetBandInfo fetches a band's info from it's ID
func GetBandInfo(bandID int) (structs.Band, error) {

	band := structs.Band{}

	members := []structs.Member{}

	mysql, err := dao.NewMysql()
	if err != nil {
		return structs.Band{}, err
	}
	defer mysql.Close()

	rows, err := mysql.Query("SELECT name, description FROM prim.bands WHERE id = ? LIMIT 1", bandID)
	if err != nil {
		return structs.Band{}, err
	}

	rows.Scan(band)

	rows, err = mysql.Query("SELECT user_id, name, role, acl FROM prim.bands_members WHERE band_id = ?", bandID)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(members[i])
		if err != nil {
			return structs.Band{}, err
		}
	}
	band.Members = members

	return band, nil
}
