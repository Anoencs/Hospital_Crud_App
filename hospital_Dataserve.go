// Data server for patientstore app
// borrows from and extends: http://www.alexedwards.net/blog/practical-persistence-sql

package main

import (
	"demo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	models.InitDB("postgres://postgres:123456@localhost:5439/demo?sslmode=disable")

	http.HandleFunc("/login", login)
	http.HandleFunc("/patients", patientsIndex)
	http.HandleFunc("/patients/show", patientsShow)
	http.HandleFunc("/patients/create", patientsCreate)
	http.HandleFunc("/patients/delete", patientsDelete)
	http.HandleFunc("/patients/search", patientSearch)

	fmt.Println("Patientstore: dataserver (port:4000)")
	http.ListenAndServe(":4000", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	user_name := r.FormValue("user_name")
	password := r.FormValue("password")

	if user_name == "" || password == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	bks := models.Verify_user(user_name, password)

	// if err != nil {
	// 	http.Error(w, http.StatusText(500), 500)
	// 	return
	// }
	b, err := json.Marshal(bks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))
}

// return an index of all Patients
func patientsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	bks, err := models.GetPatients()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//  render json
	b, err := json.Marshal(bks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))

}

// return a subset of Patient records
func patientsShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// use code as a filter
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	bks, err := models.GetPatients(code)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//  render json
	b, err := json.Marshal(bks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))

}

// create a new Patient record
func patientsCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	code := r.FormValue("code")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	addr := r.FormValue("addr")
	gender := r.FormValue("gender")
	bdate := r.FormValue("bdate")
	phonenumber := r.FormValue("phonenumber")

	if code == "" || fname == "" || lname == "" || addr == "" || gender == "" || bdate == "" || phonenumber == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	rowsAffected, err := models.CreatePatient(code, fname, lname, addr, gender, bdate, phonenumber)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// output confirmation to console
	fmt.Printf("Patient %s created successfully (%d row affected)\n", code, rowsAffected)
}

// deletes a Patient record
func patientsDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	code := r.FormValue("code")

	if code == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	rowsAffected, err := models.DeletePatient(code)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// output confirmation to console
	fmt.Printf("Patient %s deleted successfully (%d row affected)\n", code, rowsAffected)
}

// search a Patient record
func patientSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	code := r.FormValue("code")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	addr := r.FormValue("addr")
	if code == "" && fname == "" && lname == "" && addr == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	bks, err := models.SearchPatient(code, fname, lname, addr)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	b, err := json.Marshal(bks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))

}
