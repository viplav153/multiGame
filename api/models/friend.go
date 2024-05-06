package models

type User struct {
	UserID          string   `json:"userId"`
	Name            string   `json:"name"`
	OnlineStatus    string   `json:"onlineStatus"`
	Friends         []string `json:"friends"`
	SentRequest     []string `json:"sentRequest"`
	ReceivedRequest []string `json:"receivedRequest"`
}

type UserCreateRequest struct {
	Name string `json:"name"`
}

type UserUpdateRequest struct {
	OnlineStatus    string   `json:"onlineStatus"`
	Friends         []string `json:"friends"`
	SentRequest     []string `json:"sentRequest"`
	ReceivedRequest []string `json:"receivedRequest"`
}
