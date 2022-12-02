package models

import "fmt"

type User struct {
	user_name string
	password  string
}

func Verify_user(user_name, password string) (bool, error) {
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
		return false, err
	}

	defer row.Close()

	user := new(User)
	err = row.Scan(&user.user_name, &user.password)
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}
	if user.password != password {
		return false, err
	} else {
		return true, err
	}
}
