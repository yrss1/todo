package auth

type Entity struct {
	Email    *string `db:"email"`
	Password *string `db:"password"`
}
