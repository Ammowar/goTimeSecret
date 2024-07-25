package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Ammowar/goTimeSecret/storages"
)

type CreateHandler struct {
	secret_storage *storages.MemoryStorage
	templates      *template.Template
}

func NewCreateHandler(secret_storage *storages.MemoryStorage, templates *template.Template) *CreateHandler {
	return &CreateHandler{
		secret_storage: secret_storage,
		templates:      templates,
	}
}

func (ch *CreateHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	expirationStr := r.FormValue("expiration")
	var expiration time.Duration
	var err error
	if expirationStr == "" {
		expiration = 24 * time.Hour // Default to 24 hours if not specified
	} else {
		expiration, err = time.ParseDuration(expirationStr)
		if err != nil {
			http.Error(w, "Invalid expiration format", http.StatusBadRequest)
			return
		}
		if expiration > 7*24*time.Hour {
			http.Error(w, "Expiration cannot be greater than 7 days", http.StatusBadRequest)
			return
		}
	}

	id, expiresAt, err := ch.secret_storage.Create(content, expiration)
	if err != nil {
		http.Error(w, "Failed to create secret", http.StatusInternalServerError)
		return
	}

	data := struct {
		Id        string
		URL       string
		ExpiresAt string
	}{
		Id:        id,
		URL:       fmt.Sprintf("http://localhost:8080/secret/%s", id),
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}

	ch.templates.ExecuteTemplate(w, "created.html", data)
}
