package account

type UpdateProfileRequest struct {
	Name          string  `json:"name" binding:"required,min=2,max=50"`
	Email         string  `json:"email" binding:"required,email"`
	PhoneNumber   *string `json:"phone_number,omitempty" binding:"omitempty,e164"`
	ProfilePicURL *string `json:"profile_pic_url,omitempty" binding:"omitempty,url"`
	Bio           *string `json:"bio,omitempty" binding:"omitempty,max=160"`
}
