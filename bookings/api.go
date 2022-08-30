package bookings

import "net/http"

type API struct {
}

func NewAPI() *API {
	return &API{}
}

func (a API) BookFlight(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("book flight"))

}

func (a API) Bookings(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get flights"))
}
