package auth

//Login endpoint will return session ID to client
//Other endpoints will ingest session ID from client
import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/brucebales/bandmanager-backend/src/internal/structs"
)

type Session struct {
	ID         string
	UserID     int
	expiration time.Time
}

func Login(pw string, Email string, db *sql.DB) (Session, error) {
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

	return Session{
		ID:         string(sessionID[:]),
		UserID:     user.ID,
		expiration: time.Now().Add(8 * time.Hour),
	}, nil
}
