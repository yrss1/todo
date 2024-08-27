package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/yrss1/todo/internal/domain/user"
	"github.com/yrss1/todo/pkg/helpers"
	"github.com/yrss1/todo/pkg/log"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func (s *Service) Register(ctx context.Context, req user.Request) (id string, err error) {
	logger := log.LoggerFromContext(ctx).Named("Register")

	data := user.Entity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*data.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate password hash", zap.Error(err))
		return
	}
	data.Password = helpers.GetStringPtr(string(hashedPassword))

	id, err = s.userRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	return
}

func (s *Service) ValidateUser(ctx context.Context, req user.Request) (id string, err error) {
	logger := log.LoggerFromContext(ctx).Named("ValidateUser")

	data, err := s.userRepository.GetByEmail(ctx, *req.Email)
	if err != nil {
		logger.Error("failed to validate user", zap.Error(err))
		return
	}
	id = data.ID

	err = bcrypt.CompareHashAndPassword([]byte(*data.Password), []byte(*req.Password))
	if err != nil {
		logger.Error("invalid email or password", zap.Error(err))
		return
	}

	return
}

func (s *Service) GenerateJWT(ctx context.Context, id string, jwtKey []byte) (tokenString string, err error) {
	logger := log.LoggerFromContext(ctx).Named("GenerateJWT")

	expirationTime := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})

	tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		logger.Error("failed to generate token", zap.Error(err))
		return
	}

	return
}

func (s *Service) ValidateJWT(ctx context.Context, tokenString string, jwtKey []byte) (claims *Claims, err error) {
	logger := log.LoggerFromContext(ctx).Named("ValidateJWT")

	claims = &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error("invalid token signature", zap.Error(err))
			return nil, err
		}
		logger.Error("failed to parse token", zap.Error(err))
		return nil, err
	}

	if !token.Valid {
		logger.Warn("token is not valid")
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
