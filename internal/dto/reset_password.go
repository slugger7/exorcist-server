package dto

type ResetPasswordDTO struct {
	OldPassword    string `json:"oldPassword" binding:"required"`
	NewPassword    string `json:"newPassword" binding:"required"`
	RepeatPassword string `json:"repeatPassword" binding:"required,eqfield=NewPassword"`
}
