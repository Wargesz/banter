package initialisers

import "banter/models"

func SyncDatabase() {
    DB.AutoMigrate(&models.User{}, &models.Post{})
}
