package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WpsService struct {
	db *gorm.DB
}

type ExecutionStatusResponse struct {
	Status string `json:"status"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (w *WpsService) GetExecutionStatusUrl(jobId string) (string, error) {
	var url string
	result := w.db.Raw("SELECT url FROM wps_execution_status WHERE job_id = ?", jobId).Scan(&url)
	if result.Error != nil {
		return "", result.Error
	}
	return url, nil
}

func ifAllowed(url string) bool {
	allowedUrls := []string{"http://example.com", "https://example.com"}
	for _, allowedUrl := range allowedUrls {
		if url == allowedUrl {
			return true
		}
	}
	return false
}

func _getExecutionStatusResponse(url string) (*ExecutionStatusResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve execution status response")
	}

	var executionStatusResponse ExecutionStatusResponse
	err = json.NewDecoder(resp.Body).Decode(&executionStatusResponse)
	if err != nil {
		return nil, err
	}
	return &executionStatusResponse, nil
}

func _getErrorMessage(response *ExecutionStatusResponse) (*ErrorMessage, error) {
	if response.Status == "error" {
		return &ErrorMessage{Message: "Error occurred during execution"}, nil
	} else if response.Status == "failed" {
		return &ErrorMessage{Message: "Execution failed"}, nil
	} else if response.Status == "cancelled" {
		return &ErrorMessage{Message: "Execution cancelled"}, nil
	}
	return nil, nil
}

func _renderExecutionStatus(c *gin.Context, response *ExecutionStatusResponse) {
	c.JSON(http.StatusOK, response)
}

func _renderExecutionFailed(c *gin.Context, errorMessage *ErrorMessage) {
	c.JSON(http.StatusInternalServerError, errorMessage)
}

func jobReport(c *gin.Context) {
	jobId := c.Param("jobId")
	wpsService := &WpsService{db: c.MustGet("db").(*gorm.DB)}

	url, err := wpsService.GetExecutionStatusUrl(jobId)
	if err != nil {
		log.Printf("failed to retrieve execution status URL for job %s: %v", jobId, err)
		_renderExecutionFailed(c, &ErrorMessage{Message: "Failed to retrieve execution status URL"})
		return
	}

	if !ifAllowed(url) {
		log.Printf("URL %s is not allowed for job %s", url, jobId)
		_renderExecutionFailed(c, &ErrorMessage{Message: "URL is not allowed"})
		return
	}

	executionStatusResponse, err := _getExecutionStatusResponse(url)
	if err != nil {
		log.Printf("failed to retrieve execution status response for job %s: %v", jobId, err)
		_renderExecutionFailed(c, &ErrorMessage{Message: "Failed to retrieve execution status response"})
		return
	}

	errorMessage, err := _getErrorMessage(executionStatusResponse)
	if err != nil {
		log.Printf("failed to check for error in execution status response for job %s: %v", jobId, err)
		_renderExecutionFailed(c, &ErrorMessage{Message: "Failed to check for error in execution status response"})
		return
	}

	if errorMessage != nil {
		log.Printf("error occurred during execution for job %s: %s", jobId, errorMessage.Message)
		_renderExecutionFailed(c, errorMessage)
		return
	}

	_renderExecutionStatus(c, executionStatusResponse)
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	db, err := gorm.Open("postgres", "user:password@localhost/database")
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/job/:jobId/report", func(c *gin.Context) {
		c.Set("db", db)
		jobReport(c)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}