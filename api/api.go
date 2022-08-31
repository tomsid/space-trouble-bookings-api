package api

import (
	"encoding/json"
	"net/http"
	"space-trouble-bookings-api/db"
	"space-trouble-bookings-api/spacex"

	"go.uber.org/zap"
)

type API struct {
	spacex spacex.Client
	log    *zap.SugaredLogger
	db     db.Storage
}

func NewAPI(spacexClient spacex.Client, storage db.Storage, l *zap.SugaredLogger) *API {
	return &API{
		spacex: spacexClient,
		log:    l,
		db:     storage,
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (a *API) internalServerError(w http.ResponseWriter) {
	b, err := json.Marshal(ErrorResponse{Message: "Internal Server Error"})
	if err != nil {
		a.log.Error(err)
	}
	a.writeResponse(w, b)
}

func (a *API) writeJSONResponse(w http.ResponseWriter, resp interface{}) {
	b, err := json.Marshal(resp)
	if err != nil {
		a.log.Error(err)
	}
	a.writeResponse(w, b)
}

func (a *API) writeResponse(w http.ResponseWriter, resp []byte) {
	_, err := w.Write(resp)
	if err != nil {
		a.log.Error(err)
	}
}
