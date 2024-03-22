package dto

type (
	ReqAddPost struct {
		PostInHTML string   `json:"postInHtml" validate:"required,min=2,max=500"`
		Tags       []string `json:"tags" validate:"required,min=1,dive,required"`
	}
	ReqAddComment struct {
		PostID  string `json:"postId" validate:"required,uuid4"`
		Comment string `json:"comment" validate:"required,min=2,max=500"`
	}
)
