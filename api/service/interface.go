package service

import "multiGame/api/models"

type Friend interface {
	CreateUser(name string) (*models.User, error)
	AddFriend(friendID, userID string) (*models.User, error)
	RemoveFriend(friendID, userID string) (*models.User, error)
	RejectFriendRequest(friendID, userID string) (*models.User, error)
	ViewAllFriends(userID string) ([]*models.User, error)
	SendFriendRequest(friendID, userID string) (*models.User, error)
	ViewProfile(userID string) (*models.User, error)
}

type Party interface {
	CreateParty(userID string) (*models.Party, error)
	GetPartyByID(partyID string) (*models.Party, error)
	LeaveParty(userID, partyId string) (*models.Party, error)
	PartyInvitation(receiverID, partyId, verdict string) (*models.Party, error)
	SendPartyInvitation(senderId, receiverID, partyId string) (*models.Party, error)
	MakeAdminOfParty(userID, userToMakeAdmin, partyId string) (*models.Party, error)
	RemoveFromParty(userID, userToRemove, partyId string) (*models.Party, error)
}
