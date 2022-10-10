package helpers

import (
	"todomicro/database"
	"todomicro/models"
)

type User = models.User

func UserCheck(phone string, password string) (User, error) {
	row := database.DB.QueryRow("select id, name, phone, email from user where phone=? and password=? and status='active';", phone, password)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Phone, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil
}
