package request

import "devdiaries/models"

type SignupBody struct {
	User     models.User `json:"user"`
	Password string      `json:"password"`
}
