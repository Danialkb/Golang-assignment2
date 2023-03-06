package main

import (
	"assignment1/assignment2/models"
	"html/template"
	"net/http"
	"strconv"
)

var CurUser *models.User = nil

func homePage(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	if r.Method == "POST" {
		if r.FormValue("filterName") == "on" {
			fs := models.FilteringService{}
			items = fs.FilterByName()
		}
		if r.FormValue("filterPrice") == "on" {
			fs := models.FilteringService{}
			items = fs.FilterByPrice()
		}
		if r.FormValue("filterName") != "on" && r.FormValue("filterPrice") != "on" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		searching := r.URL.Query().Get("search")
		ss := models.SearchingService{}
		items = ss.Search(searching)
	}
	//fmt.Println(items)

	IsAuthorized := false
	if CurUser != nil {
		IsAuthorized = true
	}

	tmp, _ := template.ParseFiles("templates/index.html")
	tmp.Execute(w, map[string]interface{}{
		"items":        items,
		"IsAuthorized": IsAuthorized,
		"user":         CurUser,
	})
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		username := r.FormValue("username")
		password := r.FormValue("password")

		rs := models.RegistrationService{}

		rs.Register(&models.User{Name: name, Surname: surname, Password: password, Username: username})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	tmp, _ := template.ParseFiles("templates/register.html")
	tmp.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		pwd := r.FormValue("password")

		as := models.AuthorizationService{}
		var err error
		CurUser, err = as.SignIn(username, pwd)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmp, _ := template.ParseFiles("templates/login.html")
	tmp.Execute(w, nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		CurUser = nil
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func rateItemHandler(w http.ResponseWriter, r *http.Request) {
	if CurUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if r.Method == "POST" {
		rate, _ := strconv.Atoi(r.FormValue("rating"))

		rs := models.RatingService{}
		rs.RateItem(CurUser.ID, uint(id), float64(rate))

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("templates/rate.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		ID string
	}{
		ID: strconv.Itoa(id),
	}
	tmpl.Execute(w, data)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/registration", registrationHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/rate", rateItemHandler)
	http.ListenAndServe(":8000", nil)

}
