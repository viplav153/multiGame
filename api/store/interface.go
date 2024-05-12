package store

import (
	"context"
	"multiGame/api/models"
)

type FriendStore interface {
	GetUserByID(userID string) (*models.User, error)
	GetUsersFromIds(userIDs []string) ([]*models.User, error)
	CreateUser(userRequest *models.UserCreateRequest) (*models.User, error)
	AddFriend(friendUserID string, userID string) (*models.User, error)
	RemoveFriend(friendUserID string, userID string) (*models.User, error)
	RejectFriend(friendUserID string, userID string) (*models.User, error)
	SendFriendRequest(friendUserID string, userID string) (*models.User, error)
	HostNewParty(userID, partyId string) (*models.User, error)
	AddPartyInvites(userID, partyId string) (*models.User, error)
	RemovePartyInvites(userID, partyId string) (*models.User, error)
}

type RedisStore interface {
	SetKeyValue(ctx context.Context, key string, value []byte) error
	GetValue(ctx context.Context, key string) (string, error)
	IsKeyPresent(ctx context.Context, key string) bool
	SetKeyValueExpirationSame(ctx context.Context, key string, value []byte) error
}
