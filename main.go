package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

// handleFileUpload handles video file uploads
func handleFileUpload(c echo.Context) error {
	// Read form fields
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CustomResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error processing the request",
		})
	}

	files := form.File["files"]

	// Directory to store uploaded videos
	uploadDir := "./uploads/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.JSON(http.StatusInternalServerError, CustomResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error creating upload directory",
		})
	}

	var downloadLinks []string

	for _, file := range files {
		// Save the file
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, CustomResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error opening the uploaded file",
			})
		}
		defer src.Close()

		// Generate a unique filename
		filename := filepath.Join(uploadDir, file.Filename)

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
		downloadLink := fmt.Sprintf("/download/%s", file.Filename)
		downloadLinks = append(downloadLinks, downloadLink)
	}

	return c.JSON(http.StatusOK, CustomResponse{
		StatusCode: http.StatusOK,
		Message:    "Files uploaded successfully",
		Data:       downloadLinks,
	})
}

// handleDownloadLink generates download links for the uploaded videos
func handleDownloadLink(c echo.Context) error {
	filename := c.Param("filename")
	downloadLink := fmt.Sprintf("/uploads/%s", filename)
	return c.JSON(http.StatusOK, CustomResponse{
		StatusCode: http.StatusOK,
		Message:    "Download link generated successfully",
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
	e.GET("/download/:filename", handleDownloadLink)

	// Swagger and ReDoc documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/redoc", redoc.Handler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
