package access

import (
	"github.com/brucebales/bandmanager-backend/src/internal/dao"
	structs "github.com/brucebales/bandmanager-backend/src/internal/structs"
)

//GetUser accesses a user's basic info, for internal endpoint use
func GetUser(sessionId string) (structs.User, error) {

	var user structs.User

	redis := dao.NewRedis()

	mysql, err := dao.NewMysql()
	if err != nil {
		return structs.User{}, err
	}

	userid, err := redis.Get(sessionId).Result()
	if err != nil {
		return structs.User{}, err
	}

	rows, err := mysql.Query("SELECT id, name, email, password FROM prim.users WHERE id = ?", userid)

	err = rows.Scan(user)

	return user, nil
}
