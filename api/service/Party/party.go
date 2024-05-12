package Party

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"multiGame/api/models"
	"multiGame/api/service"
	"multiGame/api/store"
)

type partyService struct {
	FriendStore store.FriendStore
	RedisStore  store.RedisStore
}

func New(friendStore store.FriendStore, redisStore store.RedisStore) service.Party {
	return &partyService{FriendStore: friendStore, RedisStore: redisStore}
}

func (p *partyService) CreateParty(userID string) (*models.Party, error) {
	partyID := generatePartyID()

	_, err := p.FriendStore.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	_, err = p.FriendStore.HostNewParty(userID, partyID)
	if err != nil {
		return nil, err
	}

	// add party to redis as it is for shortest period of time
	partStruct := models.Party{Users: []models.PartyUsers{{UserID: userID, InParty: true, InviteStatus: "ACCEPTED", IsAdmin: true}}}

	// Serialize struct to JSON
	jsonData, err := json.Marshal(partStruct)
	if err != nil {
		panic(err)
	}

	err = p.RedisStore.SetKeyValue(context.Background(), partyID, jsonData)
	if err != nil {
		return nil, err
	}

	return p.GetPartyByID(partyID)
}

func (p *partyService) RemoveFromParty(userID, userToRemove, partyId string) (*models.Party, error) {
	if !p.RedisStore.IsKeyPresent(context.Background(), partyId) {
		return nil, errors.New("party done go home")
	}

	party, err := p.GetPartyByID(partyId)
	if err != nil {
		return nil, err
	}

	if !isUserPartyAdmin(party.Users, userID) {
		return nil, errors.New("bro you are not admin")
	}

	newParty := models.Party{PartyId: partyId}

	var allOtherUser []models.PartyUsers
	for _, each := range party.Users {
		if each.UserID == userToRemove {
			continue
		}

		if each.UserID == userToRemove && userToRemove == userID {
			return nil, errors.New("bro what are you doing you can't remove yourself")
		}

		allOtherUser = append(allOtherUser, each)
	}

	newParty.Users = allOtherUser
	// Serialize struct to JSON
	jsonData, err := json.Marshal(newParty)
	if err != nil {
		panic(err)
	}

	err = p.RedisStore.SetKeyValue(context.Background(), partyId, jsonData)
	if err != nil {
		return nil, err
	}

	return p.GetPartyByID(partyId)
}

func (p *partyService) MakeAdminOfParty(userID, userToMakeAdmin, partyId string) (*models.Party, error) {
	if !p.RedisStore.IsKeyPresent(context.Background(), partyId) {
		return nil, errors.New("party done go home")
	}

	party, err := p.GetPartyByID(partyId)
	if err != nil {
		return nil, err
	}

	if !isUserPartyAdmin(party.Users, userID) {
		return nil, errors.New("bro first you plead admin to make then you can make other admin")
	}

	newParty := models.Party{PartyId: partyId}

	var allOtherUser []models.PartyUsers
	for _, each := range party.Users {
		if each.UserID == userToMakeAdmin && each.IsAdmin {
			return nil, errors.New("bro what! you are already admin")
		}

		if each.UserID == userToMakeAdmin {
			each.IsAdmin = true
		}

		allOtherUser = append(allOtherUser, each)
	}

	newParty.Users = allOtherUser

	// Serialize struct to JSON
	jsonData, err := json.Marshal(newParty)
	if err != nil {
		panic(err)
	}

	err = p.RedisStore.SetKeyValue(context.Background(), partyId, jsonData)
	if err != nil {
		return nil, err
	}

	return p.GetPartyByID(partyId)
}

func (p *partyService) SendPartyInvitation(senderId, receiverID, partyId string) (*models.Party, error) {
	if !p.RedisStore.IsKeyPresent(context.Background(), partyId) {
		return nil, errors.New("party done go home")
	}

	party, err := p.GetPartyByID(partyId)
	if err != nil {
		return nil, err
	}

	isPresent := false
	for _, each := range party.Users {
		if each.UserID == senderId && each.InviteStatus == "ACCEPTED" {
			isPresent = true
			break
		}
	}

	if !isPresent {
		return nil, errors.New("bro you are not a member of this party you can't send invitation")
	}

	user, err := p.FriendStore.GetUserByID(senderId)
	if err != nil {
		return nil, err
	}

	isFriend := false
	for _, each := range user.Friends {
		if receiverID == each {
			isFriend = true
			break
		}
	}

	if !isFriend {
		return nil, errors.New("bro he is not your friend , send him friend request")
	}

	_, err = p.FriendStore.GetUserByID(receiverID)
	if err != nil {
		return nil, err
	}

	// add invite to friend profile
	_, err = p.FriendStore.AddPartyInvites(receiverID, partyId)
	if err != nil {
		return nil, err
	}

	party.Users = append(party.Users, models.PartyUsers{UserID: receiverID, InviteStatus: "PENDING"})

	// Serialize struct to JSON
	jsonData, err := json.Marshal(party)
	if err != nil {
		panic(err)
	}

	err = p.RedisStore.SetKeyValue(context.Background(), partyId, jsonData)
	if err != nil {
		return nil, err
	}

	return p.GetPartyByID(partyId)

}

func (p *partyService) PartyInvitation(receiverID, partyId, verdict string) (*models.Party, error) {
	if !p.RedisStore.IsKeyPresent(context.Background(), partyId) {
		return nil, errors.New("party done go study")
	}

	if verdict != "REJECTED" && verdict != "ACCEPTED" {
		return nil, errors.New("wrong response choose any of [ACCEPTED, REJECTED]")
	}

	party, err := p.GetPartyByID(partyId)
	if err != nil {
		return nil, err
	}

	_, err = p.FriendStore.GetUserByID(receiverID)
	if err != nil {
		return nil, err
	}

	// remove invite to friend profile
	_, err = p.FriendStore.RemovePartyInvites(receiverID, partyId)
	if err != nil {
		return nil, err
	}

	newParty := models.Party{PartyId: partyId}

	isPresent := false
	var allOtherUser []models.PartyUsers
	for _, each := range party.Users {
		if each.UserID == receiverID && each.InviteStatus == "PENDING" {
			each.InviteStatus = verdict
			isPresent = true
		}

		allOtherUser = append(allOtherUser, each)
	}

	if !isPresent {
		return nil, errors.New("bro you I guess you don't have invite for this party")
	}

	newParty.Users = allOtherUser
	// Serialize struct to JSON
	jsonData, err := json.Marshal(newParty)
	if err != nil {
		panic(err)
	}

	err = p.RedisStore.SetKeyValue(context.Background(), partyId, jsonData)
	if err != nil {
		return nil, err
	}

	return p.GetPartyByID(partyId)
}

func (p *partyService) LeaveParty(userID, partyId string) (*models.Party, error) {
	if !p.RedisStore.IsKeyPresent(context.Background(), partyId) {
		return nil, errors.New("party done go and study")
	}

	party, err := p.GetPartyByID(partyId)
	if err != nil {
		return nil, err
	}

	if isUserPartyAdmin(party.Users, userID) {
		numberOfAdmin := 0
		for _, each := range party.Users {
			if each.IsAdmin == true {
				numberOfAdmin++
			}
		}

		if numberOfAdmin == 1 {
			return nil, errors.New("bro you are only admin you can't leave the party make someone admin first")
		}
	}

	newParty := models.Party{PartyId: partyId}
	isPresent := false
	var allOtherUser []models.PartyUsers
	for _, each := range party.Users {
		if each.UserID == userID && each.InviteStatus == "ACCEPTED" {
			isPresent = true
			continue
		}

		allOtherUser = append(allOtherUser, each)
	}

	if !isPresent {
		return nil, errors.New("bro you I guess you don't have invite for this party or already left")
	}

	newParty.Users = allOtherUser
	// Serialize struct to JSON
	jsonData, err := json.Marshal(newParty)
	if err != nil {
		panic(err)
	}

	err = p.RedisStore.SetKeyValue(context.Background(), partyId, jsonData)
	if err != nil {
		return nil, err
	}

	return p.GetPartyByID(partyId)
}

// check if userIsPartAdmin
func isUserPartyAdmin(users []models.PartyUsers, userID string) bool {
	for _, each := range users {
		if each.IsAdmin && each.UserID == userID {
			return true
		}
	}

	return false
}

func (p *partyService) GetPartyByID(partyID string) (*models.Party, error) {

	party, err := p.RedisStore.GetValue(context.Background(), partyID)
	if err != nil {
		return nil, err
	}

	// Deserialize JSON back into struct
	var retrievedData models.Party
	err = json.Unmarshal([]byte(party), &retrievedData)
	if err != nil {
		panic(err)
	}

	retrievedData.PartyId = partyID

	return &retrievedData, nil
}

func generatePartyID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	partyID := make([]byte, 6)
	for i := range partyID {
		partyID[i] = charset[rand.Intn(len(charset))]
	}

	return string(partyID)
}
