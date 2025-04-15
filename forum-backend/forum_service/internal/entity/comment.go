package entity

import "time"

type comment struct {
	ID        int       `db:"id"`
	Author_id int       `db:"author_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}
