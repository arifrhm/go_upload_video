package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
)

// CustomResponse defines the response structure
type CustomResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
// @Summary Upload files
// @Description Upload video files - handles video file uploads
// @Accept multipart/form-data
// @Param files formData file true "Files to upload"
// @Success 200 {object} CustomResponse "Files uploaded successfully"
// @Router /upload [post]
func handleFileUpload(c echo.Context) error {
	// Read form fields
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CustomResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error processing the request",
		})
	}

	// Directory to store uploaded videos
	uploadDir := "./uploads/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.JSON(http.StatusInternalServerError, CustomResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error creating upload directory",
		})
	}

	// Save the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CustomResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error opening the uploaded file",
		})
	}
	defer src.Close()

	// Generate a unique filename based on current timestamp
	timestamp := time.Now().Format("20060102150405") // YYYYMMDDHHMMSS format
	ext := filepath.Ext(file.Filename) // get the file extension
	filename := filepath.Join(uploadDir, fmt.Sprintf("%s%s", timestamp, ext))

	dst, err := os.Create(filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CustomResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error creating the destination file",
		})
	}
	defer dst.Close()

	// Copy the file
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, CustomResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error copying the file",
		})
	}

	// Create download link
	downloadLink := fmt.Sprintf("/download/%s", filename)

	return c.JSON(http.StatusOK, CustomResponse{
		StatusCode: http.StatusOK,
		Message:    "File uploaded successfully",
		Data:       downloadLink,
	})
}


func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/upload", handleFileUpload)
	e.Static("/download", "./uploads")

	// Serve Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(":7070"))
}
