package services

import (
	"context"
	"errors"
	"fs-backend/config"
	"fs-backend/models"
	"fs-backend/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type ForwarderService interface {
	Signup(ctx context.Context, forwarder *models.Forwarder) error
	Login(ctx context.Context, username, password string) (string, error)
	GetForwarderDetails(ctx context.Context, username string) (*models.Forwarder, error)
	UpdateForwarderDetails(ctx context.Context, username string, forwarder *models.Forwarder) error
}

type forwarderService struct {
	repo repository.ForwarderRepository
}

func NewForwarderService(repo repository.ForwarderRepository) ForwarderService {
	return &forwarderService{repo: repo}
}

func (s *forwarderService) Signup(ctx context.Context, forwarder *models.Forwarder) error {
	// Check if user already exists
	existing, err := s.repo.FindByUsername(ctx, forwarder.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(forwarder.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	forwarder.Password = string(hashedPassword)

	// Set required empty fields securely
	forwarder.Logo = ""
	forwarder.TermsAndConditions = ""
	forwarder.DefaultLanguage = ""

	// Generate FWD ID
	fwdID, err := s.repo.GetNextForwarderID(ctx)
	if err != nil {
		return err
	}
	forwarder.ForwarderID = fwdID

	// Save to DB
	return s.repo.Create(ctx, forwarder)
}

func (s *forwarderService) Login(ctx context.Context, username, password string) (string, error) {
	forwarder, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if forwarder == nil {
		return "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(forwarder.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT Token
	secret := config.GetString("jwt.secret")
	if secret == "" {
		secret = "default_secret_key" // fallback
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"forwarderId": forwarder.ForwarderID,
		"username":    forwarder.Username,
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *forwarderService) GetForwarderDetails(ctx context.Context, username string) (*models.Forwarder, error) {
	forwarder, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if forwarder == nil {
		return nil, nil
	}
	forwarder.Password = ""
	return forwarder, nil
}

func (s *forwarderService) UpdateForwarderDetails(ctx context.Context, username string, forwarder *models.Forwarder) error {
	update := bson.M{
		"companyName":     forwarder.ForwarderCompanyName,
		"phone":           forwarder.ContactPhone,
		"email":           forwarder.Email,
		"address":         forwarder.FullAddress,
		"defaultLanguage": forwarder.DefaultLanguage,
	}
	return s.repo.UpdateByUsername(ctx, username, update)
}
