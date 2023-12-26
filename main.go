package main

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	//"fmt",
	//"log",
	//"encoding/json",
	//"image/jpeg"
	//"net/http",
	//"image/png",
	//"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Error Connecting to Database")
	e := echo.New()

	// routes
	e.POST("compress/jpeg", compressJPEG)
	e.POST("compress/png", compressPNG)

	e.Logger.Fatal(e.Start(":8080"))
}

func compressJPEG(c echo.Context) error {
	return compressImage(c, "jpeg")
}

func compressPNG(c echo.Context) error {
	return compressImage(c, "png")
}
func compressImage(c echo.Context, format string) error {

	// extract image from request
	file, err := c.FormFile("image")

	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to get image from the request.")
	}

	// open file
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to open the image.")
	}
	defer src.Close()

	// decode image
	img, _, err := image.Decode(src)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to decode the image.")
	}

	// compress the image
	var compressedImage []byte
	switch format {
	case "jpeg":
		compressedImage, err = encodeJPEG(img)
		log.Println("Compressed image is in JPEG format.")
	case "png":
		compressedImage, err = encodePNG(img)
		log.Println("Compressed image is in PNG format.")
	default:
		return c.String(http.StatusBadRequest, "Invalid image format. Supported formats: jpeg, png.")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to encode the compressed image.")
	}

	// store compressed image to given directory path
	outputPath := "/Users/shubhampatni/Desktop/workspace/image_compression/compressed-images/" + file.Filename
	err = storeImage(outputPath, compressedImage)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to store the compressed image.Failed to store the compressed image.")
	}
	// return success message
	return c.String(http.StatusOK, fmt.Sprintf("Image compressed and stored at %s", outputPath))
}

func encodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	var w = io.Writer(&buf)
	err := jpeg.Encode(w, img, &jpeg.Options{Quality: 50})
	return buf.Bytes(), err
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := png.Encode(w, img)
	return buf.Bytes(), err
}

func storeImage(outputPath string, data []byte) error {
	// Create the output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Write the compressed image to the file
	return os.WriteFile(outputPath, data, os.ModePerm)
}
