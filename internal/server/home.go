package server

import (
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/shemanaev/updtr/internal/container"
	"github.com/shemanaev/updtr/internal/templates"
)

type HomeHandler struct {
	tpl    *template.Template
	client *container.Client
}

func NewHomeHandler(client *container.Client) *HomeHandler {
	tpl, err := template.ParseFS(templates.Files, "index.gohtml", "_*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	return &HomeHandler{
		tpl,
		client,
	}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	type arguments struct {
		LastUpdate string
		Containers []container.Info
	}

	lastUpdate := h.client.GetLastUpdateTime()
	ci := h.client.GetStaleContainers()
	args := arguments{
		LastUpdate: fmt.Sprintf("%v", lastUpdate.Unix()),
		Containers: ci,
	}

	if err := h.tpl.Execute(w, args); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
