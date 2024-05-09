package friend

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"multiGame/api/models"
)

type friendStore struct {
	DB *sql.DB
}

func New(db *sql.DB) *friendStore {
	return &friendStore{DB: db}
}

func (f *friendStore) GetUserByID(userID string) (*models.User, error) {
	query := `
		SELECT userID, online_status, friends, sent_request, received_request,name
		FROM user_data
		WHERE userID = $1;
	`

	// Execute the query
	row := f.DB.QueryRow(query, userID)

	// Create a User struct to store the result
	var user models.User

	// Scan the result into the User struct
	err := row.Scan(&user.UserID, &user.OnlineStatus, pq.Array(&user.Friends), pq.Array(&user.SentRequest), pq.Array(&user.ReceivedRequest), &user.Name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &user, nil
}

func (f *friendStore) GetUsersFromIds(userIDs []string) ([]*models.User, error) {
	// Prepare the SQL query
	query := `
		SELECT userID, online_status, friends, sent_request, received_request
		FROM user_data
		WHERE userID = ANY($1);
	`

	var users []*models.User
	for _, userID := range userIDs {
		// Execute the query
		rows, err := f.DB.Query(query, pq.Array([]string{userID}))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer rows.Close()

		// Iterate over the result set
		for rows.Next() {
			// Create a User struct to store the result
			var user models.User

			// Scan the result into the User struct
			err := rows.Scan(&user.UserID, &user.OnlineStatus, pq.Array(&user.Friends), pq.Array(&user.SentRequest), pq.Array(&user.ReceivedRequest))
			if err != nil {
				log.Println(err)
				return nil, err
			}

			users = append(users, &user)
		}

		// Check for errors from iterating over rows
		if err := rows.Err(); err != nil {
			log.Println(err)
			return nil, err
		}

	}

	return users, nil
}

func (f *friendStore) CreateUser(userRequest *models.UserCreateRequest) (*models.User, error) {
	query := `INSERT INTO user_data (online_status, name) VALUES ($1, $2)RETURNING userID;`

	// Execute the SQL statement
	var lastInsertedID string
	err := f.DB.QueryRow(query, "offline", userRequest.Name).Scan(&lastInsertedID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.User{UserID: lastInsertedID, Name: userRequest.Name, OnlineStatus: "online"}, nil
}

// AddFriend accept request
func (f *friendStore) AddFriend(friendUserID string, userID string) (*models.User, error) {
	// Append a value to the friends array
	_, err := f.DB.Exec("UPDATE user_data SET received_request = array_remove(received_request, $1) WHERE userID = $2", friendUserID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Append a value to the friends array
	_, err = f.DB.Exec("UPDATE user_data SET friends = array_append(friends, $1) WHERE userID = $2", friendUserID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// add to current user to to friend's friend list also

	// Append a value to the friends array
	_, err = f.DB.Exec("UPDATE user_data SET sent_request = array_remove(sent_request, $1) WHERE userID = $2", userID, friendUserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Append a value to the friends array
	_, err = f.DB.Exec("UPDATE user_data SET friends = array_append(friends, $1) WHERE userID = $2", userID, friendUserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return f.GetUserByID(userID)
}

func (f *friendStore) RemoveFriend(friendUserID string, userID string) (*models.User, error) {

	// Append a value to the friends array
	_, err := f.DB.Exec("UPDATE user_data SET friends = array_remove(friends, $1) WHERE userID = $2", userID, friendUserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Append a value to the friends array
	_, err = f.DB.Exec("UPDATE user_data SET friends = array_remove(friends, $1) WHERE userID = $2", friendUserID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return f.GetUserByID(userID)
}

func (f *friendStore) RejectFriend(friendUserID string, userID string) (*models.User, error) {
	_, err := f.DB.Exec("UPDATE user_data SET received_request = array_remove(received_request, $1) WHERE userID = $2", friendUserID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = f.DB.Exec("UPDATE user_data SET sent_request = array_remove(sent_request, $1) WHERE userID = $2", userID, friendUserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return f.GetUserByID(userID)
}

func (f *friendStore) SendFriendRequest(friendUserID string, userID string) (*models.User, error) {
	_, err := f.DB.Exec("UPDATE user_data SET received_request = array_append(received_request, $1) WHERE userID = $2", userID, friendUserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = f.DB.Exec("UPDATE user_data SET sent_request = array_append(sent_request, $1) WHERE userID = $2", friendUserID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return f.GetUserByID(userID)
}
