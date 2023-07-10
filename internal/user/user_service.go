package user

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"server/util"
	"strconv"
	"time"
)

const (
	// TODO: migrate to env file
	secretKey = "secret"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second, // timeout for service business logic
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)

	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("Login: failed to get user for email %v: %v\n", req.Email, err.Error())
		return &LoginUserRes{}, fmt.Errorf("failed to login")
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		log.Printf("Login: password check failed for email %v \n", req.Email)
		return &LoginUserRes{}, fmt.Errorf("failed to login")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Login: failed to sign token for email %v: %v \n", req.Email, err.Error())
		return &LoginUserRes{}, fmt.Errorf("failed to login")
	}

	res := &LoginUserRes{
		accessToken: ss,
		ID:          strconv.Itoa(int(user.ID)),
		Email:       user.Email,
		Username:    user.Username,
	}

	return res, nil
}
