package service

import (
	"api-ticket/constanta"
	"api-ticket/internal/entity"
	"api-ticket/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type TalentService struct {
	talentRepository entity.ITalentRepository
}

func NewTalentService(talentRepository entity.ITalentRepository) entity.ITalentService {
	return &TalentService{
		talentRepository: talentRepository,
	}
}

func (service *TalentService) FindByID(c *gin.Context, id string) (talent entity.Talent, err error) {
	return service.talentRepository.FindByID(id)
}

func (service *TalentService) Delete(c *gin.Context, id string) (talent entity.Talent, err error) {
	talent, err = service.talentRepository.FindByID(id)
	if err != nil {
		return talent, err
	}

	// Step 2: Hapus gambar jika ada
	if talent.Photo != "" {
		oldImage := "." + talent.Photo
		if removeErr := utils.RemoveFile(oldImage); removeErr != nil {
			log.Println("Failed to delete image: ", removeErr)
			return talent, removeErr
		}
	}

	return service.talentRepository.Delete(id)
}

func (service TalentService) Create(c *gin.Context, req entity.TalentInput) (err error) {
	filename, err := utils.HandleUploadFile(c, "photo")
	if err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}

	photo := constanta.DIR_FILE + filename

	// Create
	talent := entity.Talent{
		Name:                req.Name,
		Id_promotor_created: req.Id_promotor_created,
		Photo:               photo,
		CreatedAt:           time.Now(),
	}

	if err = service.talentRepository.Create(talent); err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}

	return
}

func (service TalentService) FindAll(c *gin.Context, filter entity.RequestGetTalent) ([]entity.Talent, int64, error) {
	return service.talentRepository.FindAll(filter)
}

func (service TalentService) Update(c *gin.Context, id string, req entity.TalentInput) (talent entity.Talent, err error) {
	talentOld, err := service.talentRepository.FindByID(id)
	if err != nil {
		return talentOld, err
	}
	// Map TalentInput ke entity.Talent
	filename, err := utils.HandleUploadFile(c, "photo")
	if err != nil {
		log.Println("INI ERR ====> ", err)
		return
	}

	// Cek apakah ada file yang diunggah
	if filename != "" {
		// Hapus gambar lama jika ada
		if talentOld.Photo != "" {
			oldImage := "." + talentOld.Photo
			utils.RemoveFile(oldImage)
		}
	}

	// Jika ada file yang diunggah, set nama file yang baru
	if filename != "" {
		talent.Photo = constanta.DIR_FILE + filename
	}

	updatedBanner := entity.Talent{
		Name:                req.Name,
		Id_promotor_created: req.Id_promotor_created,
		Photo:               talent.Photo,
		UpdatedAt:           time.Now(),
	}

	return service.talentRepository.Update(id, updatedBanner)
}

// func (service TalentService) Update(c *gin.Context, id string, req entity.TalentInput) (err error) {

// }
