package models

import (
	"api/security"
	"time"
)

// User describes a user object
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:20; not null; unique" json:"nickname"`
	Email     string    `gorm:"size:50; not null; unique" json:"email"`
	Password  string    `gorm:"size:200; not null" json:"password"`
	Posts     []Post    `gorm:"foreigkey:AuthorID" json:"posts, omitempty"`
	CreatedAt time.Time `gorm:"default:current_timestamp()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}
