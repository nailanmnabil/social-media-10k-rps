package dto

type (
	ReqAddFriend struct {
		UserID string `json:"userId" validate:"required,uuid4"`
	}
	ReqDeleteFriend struct {
		UserID string `json:"userId" validate:"required,uuid4"`
	}
)
