package initialisers

import (
	"gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

var DB *gorm.DB

func DbConnect() {
    var err error
    DB, err = gorm.Open(sqlite.Open("banter.db"), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }
}
