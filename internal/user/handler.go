package user

import (
	"net/http"
)

type Handler interface {
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	GetProfileHandler(w http.ResponseWriter, r *http.Request)
	UpdateProfileHandler(w http.ResponseWriter, r *http.Request)
}
