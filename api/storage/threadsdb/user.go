package threadsdb

import (
	"context"
	"time"

	"github.com/getzion/relay/api"
	"github.com/textileio/go-threads/db"
)

func (c *threadsDbStorage) GetUsers() ([]api.User, error) {
	rawResults, err := c.Client.Find(context.Background(), *c.ID, "Users", nil, &UserSchema{})
	if err != nil {
		return nil, err
	}
	results := rawResults.([]*UserSchema)

	users := make([]api.User, len(results))
	for v, user := range results {
		userDto := api.User{
			Did:            user.Did,
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			Bio:            user.Bio,
			Img:            user.Img,
			PriceToMessage: user.PriceToMessage,
			Created:        user.Created,
			Updated:        user.Updated,
		}
		users[v] = userDto
	}

	return users, nil
}

func (c *threadsDbStorage) GetUserByDid(did string) (*api.User, error) {
	query := db.Where("did").Eq(did)
	rawResults, err := c.Client.Find(context.Background(), *c.ID, "Users", query, &UserSchema{})
	if err != nil {
		panic(err)
	}
	results := rawResults.([]*UserSchema)
	userDto := api.User{
		Did:            results[0].Did,
		Username:       results[0].Username,
		Email:          results[0].Email,
		Name:           results[0].Name,
		Bio:            results[0].Bio,
		Img:            results[0].Img,
		PriceToMessage: results[0].PriceToMessage,
		Created:        results[0].Created,
		Updated:        results[0].Updated,
	}
	return &userDto, nil
}

func (c *threadsDbStorage) InsertUser(user *api.User) error {
	currentTime := time.Now().Unix()

	userDto := &UserSchema{
		Did:            user.Did,
		Username:       user.Username,
		Email:          user.Email,
		Name:           user.Name,
		Bio:            user.Bio,
		Img:            user.Img,
		PriceToMessage: user.PriceToMessage,
		Created:        currentTime,
		Updated:        currentTime,
	}

	txn, err := c.Client.WriteTransaction(context.Background(), *c.ID, "Users")
	if err != nil {
		return err
	}
	end, err := txn.Start()
	if err != nil {
		txn.Discard()
		return err
	}

	results, err := txn.Create(user)
	if err != nil {
		txn.Discard()
		return err
	}

	// id, err := strconv.ParseInt(results[0], 10, 64)
	// if err != nil {
	// 	txn.Discard()
	// 	return err
	// }
	userDto.Id = results[0]

	err = end()
	if err != nil {
		return err
	}

	return nil
}

func (c *threadsDbStorage) GetUserByUsername(username string) (*api.User, error) {
	return nil, nil
}
