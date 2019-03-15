package structs

//Heavily tied to the auth package

type User struct {
	ID       int
	Name     string
	Email    string
	Password string //Sha256 password, comes from a DB on our end
}
