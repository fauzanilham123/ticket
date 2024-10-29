package service

import (
	"api-ticket/internal/entity"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepository entity.IAuthRepository
}

func NewAuthService(authRepository entity.IAuthRepository) entity.IAuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

func (service AuthService) Register(c *gin.Context, req entity.RegisterInputCustomer) (user entity.User, err error) {
	// Create
	Customer := entity.Customer{
		Name:      req.Name,
		Gender:    req.Gender,
		Birthday:  req.Birthday,
		CreatedAt: time.Now(),
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err // tangani error jika hashing gagal
	}

	User := entity.User{
		Name:        req.Name,
		Email:       req.Email,
		Id_type:     req.Id_type,
		Id_Promotor: req.Id_Promotor,
		Password:    string(hashedPassword),
		CreatedAt:   time.Now(),
	}

	if User, err = service.authRepository.RegisterCustomer(User, Customer); err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}
	return User, err
}

func (service AuthService) LoginCustomer(c *gin.Context, user entity.LoginInput) (entity.User, error) {
	// Proses login ke repository
	loggedInUser, err := service.authRepository.LoginCustomer(user)
	if err != nil {
		return entity.User{}, err
	}

	// Ambil waktu kedaluwarsa dan JWT secret dari environment variable
	expirationHoursStr := os.Getenv("TOKEN_HOUR_LIFESPAN")
	expirationHours, err := strconv.Atoi(expirationHoursStr)
	if err != nil {
		expirationHours = 24 // default ke 24 jam jika tidak ada atau terjadi error
	}

	jwtSecret := os.Getenv("API_SECRET")
	if jwtSecret == "" {
		return entity.User{}, errors.New("JWT_SECRET environment variable not set")
	}

	// Buat claims untuk JWT
	claims := jwt.MapClaims{
		"user_id": loggedInUser.Id,
		"exp":     time.Now().Add(time.Duration(expirationHours) * time.Hour).Unix(),
	}

	// Tanda tangan JWT dengan secret dari environment
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return entity.User{}, errors.New("failed to generate token")
	}

	// Tambahkan token ke loggedInUser
	loggedInUser.Token = tokenString

	// Return response dengan user dan token
	return loggedInUser, nil
}
