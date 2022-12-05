// Web server for patientstore app

package main

import (
	"bytes"
	"demo/models"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// dataserver url
var dataserver = "http://127.0.0.1:4000"

func main() {

	http.HandleFunc("/login", userLogin)
	http.HandleFunc("/patients", patientIndex)
	http.HandleFunc("/patients/show", patientShow)
	http.HandleFunc("/patients/create_patient", create_patient)
	http.HandleFunc("/patients/delete_patient", delete_patient)
	http.HandleFunc("/patients/search_patient", search_patient)
	http.HandleFunc("/error", patientsError)
	http.HandleFunc("/", patientsLanding)
	if ping() {
		fmt.Println("Patientstore: webserver (port:3000)")
		http.ListenAndServe(":3000", nil)
	} else {
		fmt.Printf("Dataserver not available %s \n", dataserver)
	}
}

// ping data server
func ping() bool {

	// url with endpoint
	url := fmt.Sprintf("%s", dataserver)

	_, err := http.Get(url)
	if err != nil {
		return false
	}
	return true
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html", "base.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		user_name := r.FormValue("user_name")
		password := r.FormValue("password")
		if user_name == "" || password == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		// package the data for HTTP POST
		data := url.Values{}
		data.Set("user_name", user_name)
		data.Add("password", password)

		url := fmt.Sprintf("%s/login", dataserver)
		//http_post(url, data)
		resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusFound)
		} else {
			defer resp.Body.Close()

			// read json http response
			jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			var check bool

			// unmarshal json into our struct
			err = json.Unmarshal([]byte(jsonDataFromHttp), &check)
			if err != nil {
				panic(err)
			}
			if check {
				http.Redirect(w, r, "/patients", http.StatusFound)
			} else {
				// http.Redirect(w, r, "/login", http.StatusNotFound)
				t, _ := template.ParseFiles("wrong_user.html")
				t.Execute(w, nil)
			}
		}
	}

}

// get Patients from data server
func getPatients(w http.ResponseWriter, r *http.Request, endpoint string) []models.Patient {

	// url with endpoint
	url := fmt.Sprintf("%s/%s", dataserver, endpoint)

	// GET JSON data for patients
	resp, err := http.Get(url)

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusFound)
	} else {
		defer resp.Body.Close()

		// read json http response
		jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var bks []models.Patient

		// unmarshal json into our struct
		err = json.Unmarshal([]byte(jsonDataFromHttp), &bks)
		if err != nil {
			panic(err)
		}
		return bks
	}
	return nil
}

func patientsLanding(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("welcome.html")
	t.Execute(w, nil)
}

func patientsError(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("error.html")
	t.Execute(w, nil)
}

func patientIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	// render html template
	t, _ := template.ParseFiles("index.html", "base.html")
	t.Execute(w, getPatients(w, r, "patients"))

}

func patientShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// render html template
	t, _ := template.ParseFiles("patient.html", "base.html")
	t.Execute(w, getPatients(w, r, fmt.Sprintf("patients/show?code=%s", code)))

}

// handle HTTP post to dataserver endpoint
func http_post(endpoint string, data url.Values) {

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}

}

func create_patient(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("new_patient.html", "base.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		code := r.FormValue("code")
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		addr := r.FormValue("addr")
		if code == "" || fname == "" || lname == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		// confirm price is a float
		// take the form submitted value for price

		// package the data for HTTP POST
		data := url.Values{}
		data.Set("code", code)
		data.Add("fname", fname)
		data.Add("lname", lname)
		data.Add("addr", addr)

		url := fmt.Sprintf("%s/patients/create", dataserver)
		http_post(url, data)

		http.Redirect(w, r, "/patients", http.StatusFound)
	}

}

func delete_patient(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	code := r.FormValue("code")

	if code == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// package the data for HTTP POST
	data := url.Values{}
	data.Set("code", code)

	url := fmt.Sprintf("%s/patients/delete", dataserver)
	http_post(url, data)

	http.Redirect(w, r, "/patients", http.StatusFound)

}

func search_patient(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("search_patient.html", "base.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		code := r.FormValue("code")
		fname := r.FormValue("firstname")
		lname := r.FormValue("lastname")
		addr := r.FormValue("address")
		if code == "" && fname == "" && lname == "" && addr == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		// package the data for HTTP POST
		data := url.Values{}
		data.Set("code", code)
		data.Add("firstname", fname)
		data.Add("lastname", lname)
		data.Add("address", addr)

		url := fmt.Sprintf("%s/patients/search", dataserver)
		//http_post(url, data)
		resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
		if err != nil {
			http.Redirect(w, r, "/error", http.StatusFound)
		} else {
			defer resp.Body.Close()

			// read json http response
			jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			var bks []models.Patient

			// unmarshal json into our struct
			err = json.Unmarshal([]byte(jsonDataFromHttp), &bks)
			if err != nil {
				panic(err)
			}
			t, _ := template.ParseFiles("search_res.html", "base.html")
			t.Execute(w, bks)
		}
		http.Redirect(w, r, "/patients/search_res", http.StatusFound)
	}

}
