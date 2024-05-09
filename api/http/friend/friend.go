package friend

import (
	"encoding/json"
	"multiGame/api/models"
	"multiGame/api/service"
	"net/http"
)

type friendHandler struct {
	friendService service.Friend
}

func New(friendService service.Friend) *friendHandler {
	return &friendHandler{friendService: friendService}
}

func (f *friendHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createModel models.UserCreateRequest

	err := json.NewDecoder(r.Body).Decode(&createModel)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	user, err := f.friendService.CreateUser(createModel.Name)
	getResponseOrError(w, user, err)

}

func (f *friendHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")
	friendUserId := headers.Get("friendId")

	user, err := f.friendService.AddFriend(friendUserId, userID)
	getResponseOrError(w, user, err)

}

func (f *friendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")
	friendUserId := headers.Get("friendId")

	user, err := f.friendService.RemoveFriend(friendUserId, userID)
	getResponseOrError(w, user, err)
}

func (f *friendHandler) RejectFriend(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")
	friendUserId := headers.Get("friendId")

	user, err := f.friendService.RejectFriendRequest(friendUserId, userID)
	getResponseOrError(w, user, err)
}

func (f *friendHandler) ListAllFriend(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")

	user, err := f.friendService.ViewAllFriends(userID)
	getResponseOrError(w, user, err)
}

func (f *friendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")
	friendUserId := headers.Get("friendId")

	user, err := f.friendService.SendFriendRequest(friendUserId, userID)
	getResponseOrError(w, user, err)
}

func (f *friendHandler) ViewProfile(w http.ResponseWriter, r *http.Request) {
	headers := r.Header
	userID := headers.Get("userId")

	user, err := f.friendService.ViewProfile(userID)
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
