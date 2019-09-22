package structs

type Member struct {
	UserID int //Correlates to ID in user struct
	Name   string
	Role   string
	ACL    int //Access level for editing band info, booking, etc.
}

type Band struct {
	ID          int      `db:"id"`
	Name        string   `db:"name"`
	Description string   `db:"description"`
	Members     []Member `db:"members"`
}
