package initializers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrorHandler struct {
	db  *gorm.DB
	log *log.Logger
}

func NewErrorHandler(db *gorm.DB, log *log.Logger) *ErrorHandler {
	return &ErrorHandler{db: db, log: log}
}

func (e *ErrorHandler) HandleError(c *gin.Context, err error) {
	if err == nil {
		e.log.Println("error is nil")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid error"})
		return
	}

	var gormErr *gorm.Error
	if errors.As(err, &gormErr) {
		e.log.Printf("gorm error: %v", gormErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	e.log.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (e *ErrorHandler) CustomHandleError(c *gin.Context, err error, message string) {
	if err == nil {
		e.log.Println("error is nil")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid error"})
		return
	}

	e.log.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}