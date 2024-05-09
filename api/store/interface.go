package store

import "multiGame/api/models"

type FriendStore interface {
	GetUserByID(userID string) (*models.User, error)
	GetUsersFromIds(userIDs []string) ([]*models.User, error)
	CreateUser(userRequest *models.UserCreateRequest) (*models.User, error)
	AddFriend(friendUserID string, userID string) (*models.User, error)
	RemoveFriend(friendUserID string, userID string) (*models.User, error)
	RejectFriend(friendUserID string, userID string) (*models.User, error)
	SendFriendRequest(friendUserID string, userID string) (*models.User, error)
}
