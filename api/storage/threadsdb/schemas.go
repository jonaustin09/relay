package threadsdb

import "github.com/getzion/relay/api"

type (
	UserSchema struct {
		Id             string `json:"_id,omitempty"`
		Did            string `json:"did" validate:"required"`
		Username       string `json:"username" validate:"required,username,min=6,max=16"`
		Email          string `json:"email,omitempty" validate:"omitempty,email"`
		Name           string `json:"name" validate:"required"`
		Bio            string `json:"bio,omitempty"`
		Img            string `json:"img,omitempty"`
		PriceToMessage int64  `json:"priceToMessage,omitempty"`
		Created        int64  `json:"created,omitempty"`
		Updated        int64  `json:"updated,omitempty"`
	}

	CommunitySchema struct {
		Id              string                `json:"_id,omitempty"`
		Zid             string                `json:"zid"`
		Name            string                `json:"name" validate:"required,max=150"`
		OwnerDid        string                `json:"ownerDid" validate:"required"`
		OwnerUsername   string                `json:"ownerUsername" validate:"required"`
		Description     string                `json:"description" validate:"max=250"`
		EscrowAmount    int64                 `json:"escrowAmount" validate:"gte=0,lt=100000"`
		Img             string                `json:"img,omitempty"`
		LastActive      int64                 `json:"lastActive,omitempty"`
		PricePerMessage int64                 `json:"pricePerMessage" validate:"gte=0,lt=100000"`
		PriceToJoin     int64                 `json:"priceToJoin" validate:"gte=0,lt=100000"`
		Public          api.Bool              `json:"public,omitempty"`
		Created         int64                 `json:"created,omitempty"`
		Updated         int64                 `json:"updated,omitempty"`
		Deleted         api.Bool              `json:"deleted,omitempty"`
		Tags            []string              `json:"tags,omitempty" validate:"max=5"`
		Users           []UserCommunitySchema `json:"users,omitempty"`
	}

	UserCommunitySchema struct {
		Id           string `json:"id,omitempty"`
		UserDid      string `json:"userDid"`
		CommunityZid string `json:"communityZid"`
		JoinedDate   int64  `json:"joinedDate"`
		LeftDate     int64  `json:"leftDate,omitempty"`
		LeftReason   string `json:"leftReason,omitempty"`
	}
)
