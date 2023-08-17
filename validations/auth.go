package validations

type AuthSignin struct {
	Email 		string `json:"Email" binding:"required"`
	Password 	string `json:"Password" binding:"required"`
}
