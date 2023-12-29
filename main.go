package main

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	compression "github.com/nurlantulemisov/imagecompression"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Error Connecting to Database")
	e := echo.New()

	// routes
	e.POST("compress/jpeg", compressJPEG)
	e.POST("compress/png", compressPNG)
	e.POST("compress/SVD", compressSVD)

	e.Logger.Fatal(e.Start(":8080"))
}

func compressJPEG(c echo.Context) error {
	return compressImage(c, "jpeg")
}

func compressPNG(c echo.Context) error {
	return compressImage(c, "png")
}

func compressSVD(c echo.Context) error {
	return compressWithSVD(c, "SVD")
}

func compressWithSVD(c echo.Context, format string) error {
	// extract image from request
	file, err := c.FormFile("image")
	if err != nil {
		log.Fatalf(err.Error())
	}

	// open file
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to open the image.")
	}
	defer src.Close()

	img, err := jpeg.Decode(src)

	if err != nil {
		log.Fatalf(err.Error())
	}

	compressing, _ := compression.New(50)
	compressingImage := compressing.Compress(img)

	f, err := os.Create("/Users/shubhampatni/Desktop/workspace/image_compression/svd-compressed-images/" + file.Filename)
	if err != nil {
		log.Fatalf("error creating file: %s", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(f)

	err = png.Encode(f, compressingImage)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return c.String(http.StatusOK, "image compressed and stored.")
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
	//case "SVD":
	//	//compressedImage, err = compressWithSVD(img, 20)
	//	compressing, _ := compression.New(95)
	//	compressingImage := compressing.Compress(img)
	//	log.Println("Compressed image using SVD.")

	default:
		return c.String(http.StatusBadRequest, "Invalid image format. Supported formats: jpeg, png.")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to encode the compressed image.")
	}

	// store compressed image to given directory path
	outputPath := ""
	//if format == "jpeg" || format == "png" {
	outputPath = "/Users/shubhampatni/Desktop/workspace/image_compression/compressed-images/" + file.Filename
	err = storeImage(outputPath, compressedImage)
	//} else if format == "SVD" {
	//	outputPath = "/Users/shubhampatni/Desktop/workspace/image_compression/svd-compressed-images/" + file.Filename
	//	err = storeImage(outputPath, compressingImage)
	//}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to store the compressed image.Failed to store the compressed image.")
	}
	// return success message
	return c.String(http.StatusOK, fmt.Sprintf("Image compressed and stored at %s", outputPath))
}

//func compressWithSVD(img image.Image, k int) (image.Image, error) {
//	bounds := img.Bounds()
//	rows, cols := bounds.Dy(), bounds.Dx()
//
//	// convert image to a matrix
//	matImage := mat.NewDense(rows, cols*3, nil)
//
//	for i := 0; i < rows; i++ {
//		for j := 0; j < cols; j++ {
//			r, g, b, _ := img.At(j, i).RGBA()
//			matImage.Set(i, j*3, float64(r>>8))
//			matImage.Set(i, j*3+1, float64(g>>8))
//			matImage.Set(i, j*3+2, float64(b>>8))
//		}
//	}
//
//	// Perform SVD
//	var u, v mat.Dense
//	var s mat.VecDense
//	matImage.SVD(&u, &s, &v)
//
//	// Truncate singular values and vectors
//	truncateSVD(&u, &s, &v, k)
//
//	// Reconstruct the compressed matrix
//	var compressedImage mat.Dense
//	compressedImage.Mul(&u, v.T())
//
//	// Normalize the values in the compressed matrix
//	compressedImage.Apply(func(_, _ int, v float64) float64 {
//		return math.Max(0, math.Min(255, v))
//	}, &compressedImage)
//
//	// Convert the matrix back to an image
//	resultImg := image.NewRGBA(bounds)
//	for i := 0; i < rows; i++ {
//		for j := 0; j < cols; j++ {
//			r := uint8(compressedImage.At(i, j*3))
//			g := uint8(compressedImage.At(i, j*3+1))
//			b := uint8(compressedImage.At(i, j*3+2))
//			resultImg.Set(j, i, color.RGBA{R: r, G: g, B: b, A: 255})
//		}
//	}
//
//	return resultImg, nil
//}
//
//func truncateSVD(u, s, v *mat.Dense, k int) {
//	rows, _ := u.Dims()
//
//	// Truncate singular values
//	singularValues := s.RawVector().Data
//	for i := k; i < len(singularValues); i++ {
//		singularValues[i] = 0
//	}
//
//	// Truncate rows in U and columns in V
//	u.Slice(0, rows, 0, k, u)
//	v.Slice(0, len(singularValues), 0, k, v)
//}

func encodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	var w = io.Writer(&buf)
	err := jpeg.Encode(w, img, &jpeg.Options{Quality: 10})
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
