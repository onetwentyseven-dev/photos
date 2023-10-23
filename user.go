package photos

import "time"

type User struct {
	ID        string    `structs:"id" db:"id"`
	Email     string    `structs:"email" db:"email"`
	Name      string    `structs:"name" db:"name"`
	TSCreated time.Time `structs:"ts_created" db:"ts_created"`
	TSUpdated time.Time `structs:"ts_updated" db:"ts_updated"`
}
