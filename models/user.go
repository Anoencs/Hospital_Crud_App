package models

import "fmt"

type User struct {
	User_name string
	Password  string
}

func Verify_user(user_name, password string) (bool) {
	// user := new(User)
	// var err error
	// db.QueryRow("SELECT * FROM users WHERE user_name = '$1' LIMIT 1", user_name).Scan(&user.user_name, &user.password)
	// if err != nil {
	// 	return false, err
	// }
	// defer row.Close()
	// err := row.Scan(&user.user_name, &user.password)

	query := fmt.Sprintf("SELECT * FROM users WHERE user_name = '%s' LIMIT 1", user_name)
	row, err := db.Query(query)
	if err != nil {
		return false
	}

	defer row.Close()

	user := new(User)
	row.Next()
	err = row.Scan(&user.User_name, &user.Password)
	if err != nil {
		return false
	}

	if err != nil {
		return false
	}
	if user.Password != password {
		return false
	} else {
		return true
	}
}
