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

func (f *partyHandler) CreateParty(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")

	party, err := f.partyService.CreateParty(userID)
	getResponseOrError(w, party, err)

}

func (f *partyHandler) GetPartyByID(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	partyID := headers.Get("partyId")

	party, err := f.partyService.GetPartyByID(partyID)
	getResponseOrError(w, party, err)

}

func (f *partyHandler) LeaveParty(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userID")
	partyID := headers.Get("partyId")

	party, err := f.partyService.LeaveParty(userID, partyID)
	getResponseOrError(w, party, err)

}

func (f *partyHandler) RespondToPartyInvitation(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userID")
	partyID := headers.Get("partyId")
	verdict := headers.Get("response")

	party, err := f.partyService.PartyInvitation(userID, partyID, verdict)
	getResponseOrError(w, party, err)

}

func (f *partyHandler) MakePartyAdmin(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userID")
	newAdminID := headers.Get("friendId")
	partyID := headers.Get("partyId")

	party, err := f.partyService.MakeAdminOfParty(userID, newAdminID, partyID)
	getResponseOrError(w, party, err)

}

func (f *partyHandler) KickFromParty(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userID")
	friendToKick := headers.Get("friendId")
	partyId := headers.Get("partyId")

	party, err := f.partyService.RemoveFromParty(userID, friendToKick, partyId)
	getResponseOrError(w, party, err)

}

func (f *partyHandler) SendPartyInvite(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")
	friendID := headers.Get("friendId")
	partyID := headers.Get("partyId")

	party, err := f.partyService.SendPartyInvitation(userID, friendID, partyID)
	getResponseOrError(w, party, err)
}

func getResponseOrError(w http.ResponseWriter, response interface{}, err error) {
	if err != nil {
		var errorMarshal = struct {
			Err string `json:"error"`
		}{}

		errorMarshal.Err = err.Error()
		jsonResponse, er := json.Marshal(errorMarshal)
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
