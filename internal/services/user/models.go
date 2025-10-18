package user

// User represents the database model for a user
type User struct {
	FirstName      string `json:"firstName" bson:"firstName"`
	LastName       string `json:"lastName" bson:"lastName"`
	Email          string `json:"email" bson:"email"`
	HashedPassword string `json:"hashedPassword" bson:"hashedPassword"`
	IsDeleted      bool   `json:"isDeleted" bson:"isDeleted"`
}

// SignupPayload represents the request body for user signup
type SignupPayload struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

// LoginPayload represents the request body for login
type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
