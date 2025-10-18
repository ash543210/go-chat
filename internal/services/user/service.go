package user

import (
	"context"
	"fmt"
	"time"

	"github.com/ash543210/go-chat/mongo"
	"github.com/ash543210/go-chat/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Logger *logger.Logger
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *UserService) CreateUser(firstName, lastName, email, password string, c context.Context) (User, string, error) {
	collection, err := mongo.GetCollection("users")
	if err != nil {
		return User{}, "", err
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, User{Email: email, IsDeleted: false}).Err()
	if err == nil {
		return User{}, "", fmt.Errorf("user with email %s already exists", email)
	}
	if err != mongo.ErrNoDocuments {
		return User{}, "", err // real DB error
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return User{}, "", err
	}
	document, err := collection.InsertOne(ctx, User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		HashedPassword: hashedPassword,
		IsDeleted:      false,
	})
	if err != nil {
		return User{}, "", err
	}
	s.Logger.Info(c, "Creating user: %s", email)
	user := User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		HashedPassword: hashedPassword,
		IsDeleted:      false,
	}
	return user, document.InsertedID.(primitive.ObjectID).Hex(), nil
}
