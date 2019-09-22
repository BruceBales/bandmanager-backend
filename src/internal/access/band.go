package access

import (
	"github.com/brucebales/bandmanager-backend/src/internal/dao"
	"github.com/brucebales/bandmanager-backend/src/internal/structs"
)

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
