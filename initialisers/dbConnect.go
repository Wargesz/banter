package initialisers

import (
    "gorm.io/gorm"
    "github.com/glebarez/sqlite"
)

var DB *gorm.DB

func DbConnect() {
    var err error
    DB, err = gorm.Open(sqlite.Open("banter.db"), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }
}
