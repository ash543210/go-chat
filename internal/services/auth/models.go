package auth

type RefreshToken struct {
	UserID       string `json:"userId" bson:"userId"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
	CreatedAt    int64  `json:"createdAt" bson:"createdAt"`
	ExpiresAt    int64  `json:"expiresAt" bson:"expiresAt"`
}
