package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse adalah struktur response JSON standar.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// SuccessResponse mengembalikan response sukses dengan status 200.
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse mengembalikan response sukses dengan status 201.
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengembalikan response error dengan status code tertentu.
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Success: false,
		Message: message,
	})
}

// ValidationErrorResponse mengembalikan response error validasi (422).
func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Message: "Validasi gagal",
		Errors:  errors,
	})
}

// PaginatedResponse mengembalikan response dengan pagination.
func PaginatedResponse(c *gin.Context, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "",
		Data:    data,
		Meta:    meta,
	})
}

// AppErrorResponse mengirim response berdasarkan AppError.
func AppErrorResponse(c *gin.Context, appErr *AppError) {
	if appErr.Details != nil {
		c.JSON(appErr.Code, APIResponse{
			Success: false,
			Message: appErr.Message,
			Errors:  appErr.Details,
		})
		return
	}
	c.JSON(appErr.Code, APIResponse{
		Success: false,
		Message: appErr.Message,
	})
}
