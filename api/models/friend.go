package models

type User struct {
	UserID          string   `json:"userId"`
	Name            string   `json:"name"`
	OnlineStatus    string   `json:"onlineStatus"`
	Friends         []string `json:"friends"`
	SentRequest     []string `json:"sentRequest"`
	ReceivedRequest []string `json:"receivedRequest"`
	HostedParty     []string `json:"hostedParty"`
	PartyInvites    []string `json:"partyInvites"`
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

type Party struct {
	PartyId string
	Users   []PartyUsers
}

type PartyUsers struct {
	UserID       string
	InParty      bool
	IsAdmin      bool
	InviteStatus string // invitation pending , accepted, rejected
}
