package service

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"linkbio/internal/model"
	"linkbio/internal/repo"
	"linkbio/internal/util"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidUsername    = errors.New("invalid username")
)

var usernameRegex = regexp.MustCompile(`^[a-z0-9_]{3,30}$`)

type AuthService struct {
	userRepo  *repo.UserRepo
	jwtSecret string
}

func NewAuthService(userRepo *repo.UserRepo, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*model.User, string, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, "", ErrEmailExists
	}

	hash, err := util.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user, err := s.userRepo.Create(ctx, email, hash)
	if err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !util.CheckPassword(password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *AuthService) SetUsername(ctx context.Context, userID int64, username string) (*model.User, error) {
	if !usernameRegex.MatchString(username) {
		return nil, ErrInvalidUsername
	}

	exists, err := s.userRepo.UsernameExists(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	err = s.userRepo.SetUsername(ctx, userID, username)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(ctx, userID)
}

func (s *AuthService) CheckUsername(ctx context.Context, username string) (bool, error) {
	if !usernameRegex.MatchString(username) {
		return false, nil
	}
	exists, err := s.userRepo.UsernameExists(ctx, username)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func (s *AuthService) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
