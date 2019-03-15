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

	"github.com/bbales/bandmanager-backend/src/internal/structs"
)

type Session struct {
	ID         string
	UserID     int
	expiration time.Time
}

func Login(pw string, Email string, db *sql.DB) (Session, error) {
	var user structs.User
	rows, err := db.Query("select id, name, email, password from auth.users where email = ?", Email)
	if err != nil {
		fmt.Println("Error fetching user info:", err)
		return Session{}, nil
	}
	rows.Scan(user)

	pass := sha256.Sum256([]byte(pw))

	if string(pass[:]) != user.Password {
		return Session{}, fmt.Errorf("User credentials are invalid")
	}

	seed := rand.Int()
	sessionID := sha256.Sum256([]byte(strconv.Itoa(seed)))

	return Session{
		ID:         string(sessionID[:]),
		UserID:     user.ID,
		expiration: time.Now().Add(8 * time.Hour),
	}, nil
}
