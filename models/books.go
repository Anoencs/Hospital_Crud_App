// models for patientstore app

package models

import (
	"fmt"
)

type Patient struct {
	Code  string
	Fname string
	Lname string
	Addr  string
}

// get an array of Patient with optional filter: code
func GetPatients(code ...string) ([]*Patient, error) {

	var query string
	// apply code filter if available
	if len(code) == 1 {
		query = fmt.Sprintf("SELECT * FROM patients WHERE code = '%s'", code[0])
	} else {
		query = "SELECT * FROM patients"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bks := make([]*Patient, 0)
	for rows.Next() {
		bk := new(Patient)
		err := rows.Scan(&bk.Code, &bk.Fname, &bk.Lname, &bk.Addr)
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

func CreatePatient(code string, fname string, lname string, addr string) (int64, error) {

	result, err := db.Exec("INSERT INTO patients VALUES($1, $2, $3, $4)", code, fname, lname, addr)

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

	result, err := db.Exec("DELETE FROM patients WHERE code=$1", code)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
