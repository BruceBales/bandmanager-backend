package access

import (
	"github.com/brucebales/bandmanager-backend/src/internal/dao"
	structs "github.com/brucebales/bandmanager-backend/src/internal/structs"
)

//GetUser accesses a user's basic info, for internal endpoint use
func GetUser(sessionId string) (structs.User, error) {

	var user structs.User

	redis := dao.NewRedis()
	defer redis.Close()

	mysql, err := dao.NewMysql()
	if err != nil {
		return structs.User{}, err
	}
	defer mysql.Close()

	userid, err := redis.Get(sessionId).Result()
	if err != nil {
		return structs.User{}, err
	}

	rows, err := mysql.Query("SELECT id, name, email, password FROM prim.users WHERE id = ?", userid)

	err = rows.Scan(user)

	return user, nil
}

//GetUserBands returns a list of bands that a user is a member of
//Note that this only returns the ID
func GetUserBands(userID int) ([]int, error) {
	var bands []int

	mysql, err := dao.NewMysql()
	if err != nil {
		return nil, err
	}
	defer mysql.Close()

	rows, err := mysql.Query("SELECT band_id from prim.bands_members where user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	err = rows.Scan(bands)
	if err != nil {
		return nil, err
	}

	return bands, nil
}
