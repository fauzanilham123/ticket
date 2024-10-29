package repository

import (
	"api-ticket/internal/entity"
	"errors"

	"golang.org/x/crypto/bcrypt" // Import bcrypt untuk hashingpassword
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) entity.IAuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) RegisterCustomer(user entity.User, customer entity.Customer) (entity.User, error) {
	// Mulai transaksi
	tx := repo.db.Begin()

	// Jika ada error saat memulai transaksi
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}

	// Buat customer terlebih dahulu
	if err := tx.Create(&customer).Error; err != nil {
		tx.Rollback() // rollback jika error
		return entity.User{}, err
	}

	// Mengaitkan user ke customer yang baru dibuat
	user.Id_Customer = customer.Id // Asumsi User memiliki foreign key `Id_Customer`

	// Buat user setelah customer berhasil dibuat
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback() // rollback jika error
		return entity.User{}, err
	}

	// Commit transaksi jika semua berhasil
	if err := tx.Commit().Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (repo *AuthRepository) LoginCustomer(user entity.LoginInput) (entity.User, error) {
	var User entity.User

	// Cari user berdasarkan email dan id_type = 1
	if err := repo.db.Where("email = ? AND id_type = ?", user.Email, 1).First(&User).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, errors.New("incorrect username or password")
		}
		return entity.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(user.Password)); err != nil {
		return entity.User{}, errors.New("incorrect username or password")
	}

	return User, nil
}
