package models

// User is a registered account stored in PostgreSQL.
type User struct {
	Username string `gorm:"primaryKey;size:255" json:"username"`
	Password []byte `gorm:"column:password;not null" json:"-"`
}

// UserSessions stores an active JWT session per user.
type UserSessions struct {
	Username string `gorm:"primaryKey;size:255" json:"username"`
	Token    string `gorm:"not null" json:"token"`
}

func (UserSessions) TableName() string {
	return "usersessions"
}
