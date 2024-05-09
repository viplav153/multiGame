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
