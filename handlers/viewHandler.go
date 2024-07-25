package handlers

import (
	"html/template"
	"net/http"

	"github.com/Ammowar/goTimeSecret/notifications"
	"github.com/Ammowar/goTimeSecret/storages"
)

type ViewHandler struct {
	secret_storage *storages.MemoryStorage
	templates      *template.Template
	wsConnManager  *notifications.WsConnectionManager
}

func NewViewHandler(secret_storage *storages.MemoryStorage, templates *template.Template, wsConnManager *notifications.WsConnectionManager) *ViewHandler {
	return &ViewHandler{
		secret_storage: secret_storage,
		templates:      templates,
		wsConnManager:  wsConnManager,
	}
}

func (vh *ViewHandler) Handler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/secret/"):]
	value, error := vh.secret_storage.View(id)
	vh.wsConnManager.SendMessage(id, "Viewed secret: "+id)
	if error != nil {
		vh.templates.ExecuteTemplate(w, "error.html", error.Error())
		return
	}

	vh.templates.ExecuteTemplate(w, "view.html", value)
}
