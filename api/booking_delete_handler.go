package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (a *API) BookingDelete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		a.writeBadRequest(w, ErrorResponse{Message: "booking id should be an integer and >0"})
		return
	}

	exists, err := a.db.BookingExists(ctx, id)
	if err != nil {
		a.log.Error(err)
		a.internalServerError(w)
		return
	}
	if !exists {
		a.writeBadRequest(w, ErrorResponse{Message: "booking doesn't exist"})
		return
	}

	err = a.db.BookingDelete(ctx, id)
	if err != nil {
		a.internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
