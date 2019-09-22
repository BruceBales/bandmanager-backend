package auth

//Login endpoint will return session ID to client
//Other endpoints will ingest session ID from client
import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/brucebales/bandmanager-backend/src/internal/dao"
	"github.com/brucebales/bandmanager-backend/src/internal/structs"
)

type Session struct {
	ID         string
	UserID     int
	Expiration time.Duration
}

func setSession(Session Session) {
	//Redis key is session ID, value is user ID.
	//Intention is that future functions will have
	//Session ID passed from client, and get a User ID
	//From that.
	redis := dao.NewRedis()
	redis.Set(Session.ID, Session.UserID, Session.Expiration)
}

func Login(pw string, Email string) (Session, error) {
	db, err := dao.NewMysql()
	if err != nil {
		return Session{}, err
	}

	user := structs.User{}
	rows, err := db.Query("SELECT id, name, email, password FROM prim.users WHERE email = ? limit 1", Email)
	if err != nil {
		fmt.Println("Error fetching user info:", err)
		return Session{}, nil
	}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			fmt.Println("Error scanning user data: ", err)
		}
	}

	sha256_hash := fmt.Sprintf("%x", sha256.Sum256([]byte(pw)))
	if string(sha256_hash) != user.Password {
		return Session{}, fmt.Errorf("User credentials are invalid")
	}

	seed := rand.Int()
	sessionID := fmt.Sprintf("%x", sha256.Sum256([]byte(strconv.Itoa(seed))))

	session := Session{
		ID:         string(sessionID[:]),
		UserID:     user.ID,
		Expiration: 8 * time.Hour,
	}
	//Sending session to Redis
	setSession(session)
	return session, nil
}

func CreateUser(name, email, password string) error {
	db, err := dao.NewMysql()
	if err != nil {
		return err
	}
	// Sha256-ing password before insert
	pw := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	_, err = db.Exec("INSERT INTO prim.users (name, email, password) VALUES (?, ?, ?)", name, email, pw)
	if err != nil {
		return err
	}
	return nil
}
