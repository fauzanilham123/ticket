package controllers

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseController struct {
}

type Response struct {
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
	Pagination Paginations `json:"pagination"`
}

type Paginations struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	LastPage    int   `json:"last_page"`
	NextPage    int   `json:"next_page"`
	PrevPage    int   `json:"prev_page"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func SendResponse(ctx *gin.Context, result interface{}, message string) {
	response := Response{
		Success: true,
		Code:    200,
		Message: message,
		Result:  result,
	}

	ctx.JSON(http.StatusOK, response)
}
func SendResponseGetAll(ctx *gin.Context, result interface{}, message string, pagination Paginations) {
	response := Response{
		Success:    true,
		Code:       200,
		Message:    message,
		Result:     result,
		Pagination: pagination,
	}

	ctx.JSON(http.StatusOK, response)
}

func SendError(ctx *gin.Context, errorMessage string, errorMessages string, code int) {
	response := Response{
		Success: false,
		Code:    code,
		Message: errorMessage,
		Result:  nil,
	}

	if len(errorMessages) > 0 {
		response.Result = errorMessages
	}

	ctx.JSON(code, response)
}

func ExtractPagination(c *gin.Context) Pagination {
	page := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("perPage", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	perPageInt, err := strconv.Atoi(perPage)
	if err != nil || perPageInt < 1 {
		perPageInt = 10
	}

	return Pagination{Page: pageInt, PerPage: perPageInt}
}

func PaginateQuery(db *gorm.DB, pagination Pagination) *gorm.DB {
	return db.Offset((pagination.Page - 1) * pagination.PerPage).Limit(pagination.PerPage)
}

func GetColumns(model interface{}) map[string]bool {
	validColumns := make(map[string]bool)

	// Mengambil tipe dari model (misalnya models.Banner)
	t := reflect.TypeOf(model)

	// Loop melalui semua field dalam struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Mengambil nama field dan menambahkannya ke map validColumns
		jsonTag := field.Tag.Get("json")
		// Mengambil nama field JSON (atau gunakan nama field langsung jika tidak ada tag json)
		if jsonTag != "" && jsonTag != "-" {
			validColumns[strings.Split(jsonTag, ",")[0]] = true
		} else {
			validColumns[field.Name] = true
		}
	}

	return validColumns
}
