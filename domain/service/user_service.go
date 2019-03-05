package service

import (
	"context"
	irepo "github.com/tozastation/gRPC-Training-Golang/domain/repository"
	"github.com/tozastation/gRPC-Training-Golang/infrastructure/persistence/model/db"
	"github.com/tozastation/gRPC-Training-Golang/interfaces/auth"
	rpc_user "github.com/tozastation/gRPC-Training-Golang/interfaces/rpc/user"
)

// IUserService ...
type IUserService interface {
	GetMe(ctx context.Context, token string) (*rpc_user.GetUser, error)
	SignIn(ctx context.Context, uID, password string) (string, error)
	SignUp(ctx context.Context, user *rpc_user.PostUser) (string, error)
}

type userService struct {
	irepo.IUserRepository
}

// NewUserService is ...
func NewUserService(repo irepo.IUserRepository) IUserService {
	return &userService{repo}
}

func (srv *userService) GetMe(ctx context.Context, token string) (*rpc_user.GetUser, error) {
	user, err := srv.IUserRepository.FindUserByUserToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return dbToPostUser(user), nil
}

func (srv *userService) SignUp(ctx context.Context, user *rpc_user.PostUser) (string, error) {
	dbUser, err := postUserToDB(user)
	if err != nil {
		return "", err
	}
	token, err := srv.IUserRepository.CreateUser(dbUser)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (srv *userService) SignIn(ctx context.Context, uID, password string) (string, error) {
	hashedPass, err := auth.Hashed(password)
	if err != nil {
		return "", err
	}
	token, err := srv.IUserRepository.Login(uID, hashedPass)
	if err != nil {
		return "", err
	}
	return token, nil
}

func postUserToDB(user *rpc_user.PostUser) (*db.User, error) {
	password, err := auth.Hashed(user.GetPassword())
	if err != nil {
		return nil, err
	}
	return &db.User{
		Name:        user.GetName(),
		CityName:    user.GetCityName(),
		Password:    password,
		AccessToken: auth.CreateJWT(user.GetName(), user.GetCityName()),
	}, nil
}

func dbToPostUser(user *db.User) *rpc_user.GetUser {
	return &rpc_user.GetUser{
		Name:     user.Name,
		CityName: user.CityName,
	}
}
