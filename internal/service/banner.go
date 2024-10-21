package service

import (
	"api-ticket/constanta"
	"api-ticket/internal/entity"
	"api-ticket/utils"
	"errors"
	"log"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type BannerService struct {
	bannerRepository entity.IBannerRepository
}

func NewBannerService(bannerRepository entity.IBannerRepository) entity.IBannerService {
	return &BannerService{
		bannerRepository: bannerRepository,
	}
}

func handleUploadFile(c *gin.Context, form string) (string, error) {
	fileHeader, _ := c.FormFile(form)
	var filename string

	// Cek apakah ada file yang diunggah
	if fileHeader != nil {
		fileExtention := []string{".png", ".jpg", ".jpeg", ".svg"}
		isFileValidated := utils.FileValidationByExtension(fileHeader, fileExtention)
		if !isFileValidated {
			return "", errors.New("file not allowed")
		}

		// Batas ukuran file dalam byte (2MB)
		maxFileSize := int64(2 * 1024 * 1024) // 2MB

		// Periksa ukuran file
		if fileHeader.Size > maxFileSize {
			return "", errors.New("File size too large (max 2MB)")
		}

		extensionFile := filepath.Ext(fileHeader.Filename)
		filename = utils.RandomFileName(extensionFile)

		isSaved := utils.SaveFile(c, fileHeader, filename)
		if !isSaved {
			return "", errors.New("Internal server error, can't save file")
		}
	}

	// Jika tidak ada file yang diunggah, filename akan tetap kosong

	return filename, nil
}

func (service BannerService) Create(c *gin.Context, req entity.BannerInput) (err error) {
	filename, err := handleUploadFile(c, "img")
	if err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}

	img := constanta.DIR_FILE + filename

	// Create
	banner := entity.Banner{
		Title:     req.Title,
		Slug:      req.Slug,
		Desc:      req.Desc,
		Img:       img,
		CreatedAt: time.Now(),
	}

	if err = service.bannerRepository.Create(banner); err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}

	return
}

func (service BannerService) FindAll(c *gin.Context, filter entity.RequestGetBanner) ([]entity.Banner, int64, error) {
	return service.bannerRepository.FindAll(filter)
}
