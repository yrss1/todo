package task

type Entity struct {
	ID          string  `db:"id"`
	UserID      *string `db:"user_id"`
	Title       *string `db:"title"`
	Description *string `db:"description"`
	Status      *string `db:"status"`
	DueDate     *string `db:"due_date"`
}
