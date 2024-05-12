package party

import (
	"encoding/json"
	"multiGame/api/service"
	"net/http"
)

type partyHandler struct {
	partyService service.Party
}

func New(partyService service.Party) *partyHandler {
	return &partyHandler{partyService: partyService}
}

func (f *partyHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")

	user, err := f.partyService.CreateParty(userID)
	getResponseOrError(w, user, err)

}

func getResponseOrError(w http.ResponseWriter, response interface{}, err error) {
	if err != nil {
		jsonResponse, er := json.Marshal(err)
		if er != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
	} else {
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}
