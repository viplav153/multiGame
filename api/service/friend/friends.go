package friend

import (
	"errors"
	"multiGame/api/models"
	"multiGame/api/service"
	"multiGame/api/store"
)

type friendService struct {
	FriendStore store.FriendStore
}

func New(friendStore store.FriendStore) service.Friend {
	return &friendService{FriendStore: friendStore}
}

func (f *friendService) CreateUser(name string) (*models.User, error) {
	return f.FriendStore.CreateUser(&models.UserCreateRequest{Name: name})
}

func (f *friendService) AddFriend(friendID, userID string) (*models.User, error) {
	friend, err := f.FriendStore.GetUserByID(friendID)
	if err != nil {
		return nil, err
	}

	if friend == nil {
		return nil, errors.New("friend Id not exists")
	}

	isPresent := false
	for _, each := range friend.SentRequest {
		if each == userID {
			isPresent = true
			break
		}
	}

	if isPresent == false {
		return nil, errors.New("not a authorise action")
	}

	return f.FriendStore.AddFriend(friendID, userID)
}

func (f *friendService) RemoveFriend(friendID, userID string) (*models.User, error) {
	friend, err := f.FriendStore.GetUserByID(friendID)
	if err != nil {
		return nil, err
	}

	if friend == nil {
		return nil, errors.New("friend Id not exists")
	}

	isPresent := false
	for _, each := range friend.Friends {
		if each == userID {
			isPresent = true
			break
		}
	}

	if isPresent == false {
		return nil, errors.New("you are already not friend")
	}

	return f.FriendStore.RemoveFriend(friendID, userID)
}

func (f *friendService) RejectFriendRequest(friendID, userID string) (*models.User, error) {
	friend, err := f.FriendStore.GetUserByID(friendID)
	if err != nil {
		return nil, err
	}

	if friend == nil {
		return nil, errors.New("friend Id not exists")
	}

	isPresent := false
	for _, each := range friend.SentRequest {
		if each == userID {
			isPresent = true
			break
		}
	}

	if isPresent == false {
		return nil, errors.New("you didn't get request from this ID")
	}

	return f.FriendStore.RejectFriend(friendID, userID)
}

func (f *friendService) SendFriendRequest(friendID, userID string) (*models.User, error) {
	friend, err := f.FriendStore.GetUserByID(friendID)
	if err != nil {
		return nil, err
	}

	if friend == nil {
		return nil, errors.New("friend Id not exists")
	}

	for _, each := range friend.ReceivedRequest {
		if each == userID {
			return nil, errors.New("friend request already send")
		}
	}

	return f.FriendStore.SendFriendRequest(friendID, userID)
}

func (f *friendService) ViewAllFriends(userID string) ([]*models.User, error) {
	user, err := f.FriendStore.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return f.FriendStore.GetUsersFromIds(user.Friends)
}

func (f *friendService) ViewProfile(userID string) (*models.User, error) {
	return f.FriendStore.GetUserByID(userID)
}
