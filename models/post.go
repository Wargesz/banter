package models

import (
    "gorm.io/gorm"
)

type Post struct {
    gorm.Model
    Title string
    Content string
    Upvotes uint
    Downvotes uint
    User User
    UserID uint
}
