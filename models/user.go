package models

import "fmt"

type User struct {
	user_name string
	password  string
}

func Verify_user(user_name, password string) (bool, error) {
	row, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE user_name = '%s'", user_name))
	if err != nil {
		return false, err
	}
	user := new(User)
	err = row.Scan(&user.user_name, &user.password)
	if err != nil {
		return false, err
	}
	if user.password != password {
		return false, err
	} else {
		return true, err
	}
}
