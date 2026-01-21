// Package models defines domain entities persisted by GORM and returned by the API.
package models

// User is a registered account. Password is stored as a bcrypt hash and never
// serialized to JSON (json:"-").
type User struct {
	Username string `gorm:"primaryKey;size:255" json:"username"`
	Password []byte `gorm:"column:password;not null" json:"-"`
}

// UserSessions stores one active JWT string per username so logout can revoke access.
type UserSessions struct {
	Username string `gorm:"primaryKey;size:255" json:"username"`
	Token    string `gorm:"not null" json:"token"`
}

// TableName maps UserSessions to the "usersessions" table (matches init.sql).
func (UserSessions) TableName() string {
	return "usersessions"
}
