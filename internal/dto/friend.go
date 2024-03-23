package dto

type (
	ReqAddFriend struct {
		UserID string `json:"userId" validate:"required,uuid4"`
	}
	ReqDeleteFriend struct {
		UserID string `json:"userId" validate:"required,uuid4"`
	}
	ParamGetFriends struct {
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
		SortBy     string `json:"sortBy" validate:"oneof=friendCount createdAt"`
		OrderBy    string `json:"orderBy" validate:"oneof=asc desc"`
		OnlyFriend bool   `json:"onlyFriend"`
		Search     string `json:"search"`
	}
	ResGetFriends struct {
		UserID      string `json:"userId"`
		Name        string `json:"name"`
		ImageURL    string `json:"imageUrl"`
		FriendCount int    `json:"friendCount"`
		CreatedAt   string `json:"createdAt"`
	}
)
