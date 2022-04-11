package threadsdb

import (
	"context"
	"time"

	"github.com/getzion/relay/api"
)

func (c *threadsDbStorage) GetCommunities() ([]api.Community, error) {

	rawResults, err := c.Client.Find(context.Background(), *c.ID, "Communities", nil, &CommunitySchema{})
	if err != nil {
		return nil, err
	}
	results := rawResults.([]*CommunitySchema)

	communities := make([]api.Community, len(results))
	for v, community := range results {
		communityDto := api.Community{
			Zid:             community.Zid,
			Name:            community.Name,
			OwnerDid:        community.OwnerDid,
			OwnerUsername:   community.OwnerUsername,
			Description:     community.Description,
			EscrowAmount:    community.EscrowAmount,
			LastActive:      community.LastActive,
			PricePerMessage: community.PricePerMessage,
			PriceToJoin:     community.PriceToJoin,
			Public:          community.Public,
			Deleted:         community.Deleted,
			Tags:            community.Tags,
			Img:             community.Img,
			Created:         community.Created,
			Updated:         community.Updated,
		}
		communities[v] = communityDto
	}

	return communities, nil
}

func (c *threadsDbStorage) GetCommunityByZid(string) (*api.Community, error) {
	// var community v1.CommunityORM
	// result := s.connection.DB.Model(&community).First(&community, "zid = ?", zid)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	// return &community, nil
	return nil, nil
}

func (c *threadsDbStorage) InsertCommunity(community *api.Community) error {
	currentTime := time.Now().Unix()

	communityDto := &CommunitySchema{
		Zid:             community.Zid,
		Name:            community.Name,
		OwnerDid:        community.OwnerDid,
		OwnerUsername:   community.OwnerUsername,
		Description:     community.Description,
		EscrowAmount:    community.EscrowAmount,
		LastActive:      community.LastActive,
		PricePerMessage: community.PricePerMessage,
		PriceToJoin:     community.PriceToJoin,
		Public:          community.Public,
		Deleted:         community.Deleted,
		Tags:            community.Tags,
		Img:             community.Img,
		Created:         currentTime,
		Updated:         currentTime,
	}

	txn, err := c.Client.WriteTransaction(context.Background(), *c.ID, "Communities")
	if err != nil {
		return err
	}
	end, err := txn.Start()
	if err != nil {
		txn.Discard()
		return err
	}

	results, err := txn.Create(community)
	if err != nil {
		txn.Discard()
		return err
	}

	// id, err := strconv.ParseInt(results[0], 10, 64)
	// if err != nil {
	// 	txn.Discard()
	// 	return err
	// }
	communityDto.Id = results[0]

	err = end()
	if err != nil {
		return err
	}

	return nil
}

func (c *threadsDbStorage) AddUserToCommunity(communityZid, userDid string) error {
	// association := s.connection.DB.Model(community).Omit("users").Association("users").Append(user)
	// return association.Error
	return nil
}

func (c *threadsDbStorage) RemoveUserToCommunity(communityZid, userDid, leftReason string) error {
	// association := s.connection.DB.Model(community).Omit("users").Association("users").Delete(user)
	// return association.Error
	return nil
}
