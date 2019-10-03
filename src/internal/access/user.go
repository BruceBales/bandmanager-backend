package access

import (
	"database/sql"
	"fmt"

	"github.com/brucebales/bandmanager-backend/src/internal/dao"
	structs "github.com/brucebales/bandmanager-backend/src/internal/structs"
)

//GetUser accesses a user's basic info, for internal endpoint use
func GetUser(sessionId string, db *sql.DB) (structs.User, error) {

	var user = structs.User{}

	redis := dao.NewRedis()
	defer redis.Close()

	userid, err := redis.Get(sessionId).Result()
	if err != nil {
		return structs.User{}, err
	}

	rows, err := db.Query("SELECT id, name, email FROM prim.users WHERE id = ?", userid)

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			fmt.Println(err)
		}
	}

	return user, nil
}

func GetUserByID(userId int, db *sql.DB) (structs.User, error) {
	var user = structs.User{}

	rows, err := db.Query("SELECT id, name, email FROM prim.users WHERE id = ?", userId)

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			fmt.Println(err)
		}
	}

	return user, nil

}

//GetUserBands returns a list of bands that a user is a member of
//Note that this only returns the ID
func GetUserBands(userID int, db *sql.DB) ([]int, error) {
	var bands []int

	rows, err := db.Query("SELECT band_id from prim.bands_members where user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	err = rows.Scan(bands)
	if err != nil {
		return nil, err
	}

	return bands, nil
}
