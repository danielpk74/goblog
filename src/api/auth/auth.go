package auth

import (
	"api/database"
	"api/models"
	"api/security"
	"api/utils/channels"

	"github.com/jinzhu/gorm"
)

func SignIn(email, password string) (string, error) {
	user := models.User{}
	var err error
	done := make(chan bool)
	var db *gorm.DB

	go func(ch chan<- bool) {
		defer close(ch)
		db, err = database.Connect()
		if err != nil {
			ch <- false
			return
		}
		defer db.Close()

		err = db.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
		if err != nil {
			ch <- false
			return
		}

		err = security.VerifyPassword(user.Password, password)
		if err != nil {
			ch <- false
			return
		}

		ch <- true
	}(done)

	if channels.OK(done) {
		return CreateToken(user.ID)
	}

	return "", err
}
