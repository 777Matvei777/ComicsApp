package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"time"
)

type LoginPageData struct {
	Error string
}

type ComicsPageData struct {
	SearchQuery string
	Comics      []string
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ComicsURLs struct {
	Urls []string `json:"urls"`
}

func main() {
	templates := template.Must(template.ParseGlob("web-server/templates/*.html"))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data := LoginPageData{}
			templates.ExecuteTemplate(w, "login.html", data)
		case http.MethodPost:

			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusInternalServerError)
				return
			}
			username := r.FormValue("username")
			password := r.FormValue("password")

			loginData := LoginBody{
				Username: username,
				Password: password,
			}

			loginDataJSON, err := json.Marshal(loginData)
			if err != nil {
				http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
				return
			}

			resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(loginDataJSON))
			if err != nil {
				http.Error(w, "Error POST request", http.StatusInternalServerError)
				return
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "Error reading response body", http.StatusInternalServerError)
				return
			}
			tokenString := string(body)

			http.SetCookie(w, &http.Cookie{
				Name:    "Authorization",
				Value:   tokenString,
				Expires: time.Now().Add(12 * time.Hour),
			})

			http.Redirect(w, r, "/comics", http.StatusFound)
		}
	})

	http.HandleFunc("/comics", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusInternalServerError)
				return
			}
			search := r.FormValue("search")
			searchQuery := fmt.Sprintf("http://localhost:8080/pics?search=%s", url.QueryEscape(search))
			resp, err := http.Get(searchQuery)
			if err != nil {
				http.Error(w, "Error GET comics", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "Error parse comics", http.StatusInternalServerError)
				return
			}

			var Comics ComicsURLs
			err = json.Unmarshal(body, &Comics)
			if err != nil {
				http.Error(w, "Error Unmarshaling comics", http.StatusInternalServerError)
				return
			}

			data := ComicsPageData{
				SearchQuery: search,
				Comics:      Comics.Urls,
			}
			templates.ExecuteTemplate(w, "comics.html", data)
		}
	})

	http.ListenAndServe(":8081", nil)
}
