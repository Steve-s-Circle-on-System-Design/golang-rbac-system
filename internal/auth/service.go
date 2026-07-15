package auth

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Steve-s-Circle-on-System-Design/golang-rbac-system/internal/users"
	usersdb "github.com/Steve-s-Circle-on-System-Design/golang-rbac-system/internal/users/sqlc"
)

var (
	ErrUserWithEmailAlreadyExists  = errors.New("user with that email already exists")
	ErrNonExistentUser             = errors.New("user doesn't exist with that email")
	ErrPasswordMismatchDuringLogin = errors.New("invalid password")
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Service interface {
	RegisterWithPassword(ctx context.Context, email, password string) error
	LoginWithPassword(ctx context.Context, email, password string) error
}

type authService struct {
	userRepository *users.Repository
}

func NewService(userRepository *users.Repository) Service {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) RegisterWithPassword(ctx context.Context, email, password string) error {
	_, err := s.userRepository.GetByEmail(ctx, email)
	if err == nil {
		return ErrUserWithEmailAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		log.Println("failed to check existing user:", err)
		return err
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Something went wrong while hashing the password", err.Error())
		return err
	}

	_, err = s.userRepository.Create(ctx, usersdb.CreateUserParams{
		Email:        email,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		log.Println("Something went wrong while trying to save the new user in the db", err.Error())
		return err
	}
	return nil
}

func (s *authService) LoginWithPassword(ctx context.Context, email, password string) error {
	existingUser, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNonExistentUser
		}
		log.Println("failed to check existing user:", err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password))
	if err != nil {
		return ErrPasswordMismatchDuringLogin
	}
	return nil
}
