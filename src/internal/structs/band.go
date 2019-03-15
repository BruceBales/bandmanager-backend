package structs

type member struct {
	userID int //Correlates to ID in user struct
	name   string
	role   string
}

type Band struct {
	ID          int
	name        string
	description string
	members     []member
}
