package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique"`
    Password string
    ProfilePicture string
    Posts []Post
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
    tx.Clauses(clause.Returning{}).Where("user_id == ?", u.ID).Delete(&Post{})
    return
}
