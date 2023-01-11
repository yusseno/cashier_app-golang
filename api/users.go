package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	// Read username and password request with FormValue.
	creds := model.Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	} // TODO: replace this

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
	}
	// TODO: answer here

	err := api.usersRepo.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	filepath := path.Join("views", "status.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var data = map[string]string{"name": creds.Username, "message": "register success!"}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	// Read usernmae and password request with FormValue.
	creds := model.Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	} // TODO: replace this

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
	}
	// TODO: answer here

	err := api.usersRepo.LoginValid(creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	// Generate Cookie with Name "session_token", Path "/", Value "uuid generated with github.com/google/uuid", Expires time to 5 Hour.
	// TODO: answer here
	cookie := uuid.New().String()

	session := model.Session{
		Token:    cookie,
		Username: r.FormValue("username"),
		Expiry:   time.Now().Add(5 * time.Hour),
	} // TODO: replace this

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   cookie,
		Path:    "/",
		Expires: session.Expiry,
	})
	err = api.sessionsRepo.AddSessions(session)

	api.dashboardView(w, r)
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	//Read session_token and get Value:
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value // TODO: replace this
	fmt.Println(sessionToken)

	api.sessionsRepo.DeleteSessions(sessionToken)
	c.Name = "session_token"
	c.Expires = time.Unix(0, 0)
	c.MaxAge = -1
	http.SetCookie(w, c)

	//Set Cookie name session_token value to empty and set expires time to Now:
	// TODO: answer here

	filepath := path.Join("views", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}
