package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shemanaev/updtr/internal/container"
)

type ApiHandler struct {
	client *container.Client
}

func NewApiHandler(client *container.Client) *ApiHandler {
	return &ApiHandler{
		client: client,
	}
}

func (h *ApiHandler) ListContainers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		LastUpdate string
		Containers []container.Info
	}

	afterQuery := r.URL.Query().Get("after")
	if afterQuery == "" {
		afterQuery = "0"
	}

	afterInt, err := strconv.ParseInt(afterQuery, 10, 64)
	if err != nil {
		afterInt = 0
	}
	after := time.Unix(afterInt, 0).Local()

	for i := 0; i < 10; i++ {
		lastUpdate := h.client.GetLastUpdateTime()

		if lastUpdate.Sub(after).Seconds() > 1 {
			ci := h.client.GetStaleContainers()
			res := response{
				LastUpdate: fmt.Sprintf("%v", lastUpdate.Unix()),
				Containers: ci,
			}

			if err := JSON(w, res); err != nil {
				InternalError(w)
			}
			return
		}

		time.Sleep(time.Second)
	}

	NotModified(w)
}

func (h *ApiHandler) UpdateContainers(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Ids []string
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, err.Error())
		return
	}

	h.client.Update(req.Ids)

	Accepted(w)
}

func (h *ApiHandler) RefreshContainers(w http.ResponseWriter, r *http.Request) {
	go h.client.RefreshContainers()

	Accepted(w)
}
