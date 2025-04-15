package entity

type Post struct {
	ID       int    `json:"id" db:"id"`
	AuthorId int    `json:"author_id" db:"author_id"`
	Title    string `json:"title" db:"title"`
	Content  string `json:"content" db:"content"`
}
