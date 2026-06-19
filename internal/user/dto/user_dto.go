package dto

type GetUserRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
