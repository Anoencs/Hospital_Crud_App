// models for patientstore app

package models

import (
	"fmt"
	"log"
)

type Patient struct {
	Code        string
	Firstname   string
	Lastname    string
	Address     string
	Gender      string
	Bdate       string
	Phonenumber string
}

// get an array of Patient with optional filter: code
func GetPatients(code ...string) ([]*Patient, error) {

	var query string
	// apply code filter if available
	if len(code) == 1 {
		query = fmt.Sprintf("SELECT * FROM patient WHERE code = '%s'", code[0])
	} else {
		query = "SELECT * FROM patient"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bks := make([]*Patient, 0)
	for rows.Next() {
		bk := new(Patient)
		err := rows.Scan(&bk.Code, &bk.Firstname, &bk.Lastname, &bk.Address,&bk.Gender, &bk.Bdate, &bk.Phonenumber)
		if err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bks, nil
}

func CreatePatient(code, fname, lname, addr , gender, bdate, phonenumber string) (int64, error) {

	result, err := db.Exec("INSERT INTO patient VALUES($1, $2, $3, $4, $5, $6, $7)", code, fname, lname, addr, gender, bdate, phonenumber)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func DeletePatient(code string) (int64, error) {

	result, err := db.Exec("DELETE FROM patient WHERE code=$1", code)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func SearchPatient(code string, fname string, lname string, addr string) ([]*Patient, error) {
	query := ""
	sql := "SELECT * FROM patient WHERE"

	if code != "" {
		query = fmt.Sprintf("%s code = '%s' AND", query, code)
	}
	if fname != "" {
		query = fmt.Sprintf("%s firstname = '%s' AND", query, fname)
	}
	if lname != "" {
		query = fmt.Sprintf("%s lastname = '%s' AND", query, lname)
	}
	if addr != "" {
		query = fmt.Sprintf("%s address = '%s' AND", query, addr)
	}

	if query != "" {
		sql = sql + query
	}

	sql = sql[0 : len(sql)-4]

	result, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()
	bks := make([]*Patient, 0)
	for result.Next() {
		bk := new(Patient)
		err := result.Scan(&bk.Code, &bk.Firstname, &bk.Lastname, &bk.Address,&bk.Gender, &bk.Bdate, &bk.Phonenumber)
		if err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = result.Err(); err != nil {
		return nil, err
	}

	if err != nil {
		return []*Patient{}, err
	}

	// rowsAffected, err := result.RowsAffected()

	if err != nil {
		return []*Patient{}, err
	}

	return bks, nil
}
