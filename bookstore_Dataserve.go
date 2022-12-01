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
	models.InitDB("postgres://postgres:1@localhost:5439/demo?sslmode=disable")

	http.HandleFunc("/patients", patientsIndex)
	http.HandleFunc("/patients/show", patientsShow)
	http.HandleFunc("/patients/create", patientsCreate)
	http.HandleFunc("/patients/delete", patientsDelete)

	fmt.Println("Patientstore: dataserver (port:4000)")
	http.ListenAndServe(":4000", nil)
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
	if code == "" || fname == "" || lname == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	rowsAffected, err := models.CreatePatient(code, fname, lname, addr)

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
