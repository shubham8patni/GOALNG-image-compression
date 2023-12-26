<br>
<b> use the image package in Golang to determine the format of an image. 



    package main
    
    import (
    "fmt"
    "image"
    _ "image/jpeg" // Importing for JPEG support
    _ "image/png"  // Importing for PNG support
    "os"
    )
    
    func main() {
    filePath := "image.png" // Replace with the actual file path
    
        // Open the image file
        file, err := os.Open(filePath)
        if err != nil {
            fmt.Println("Error opening image file:", err)
            return
        }
        defer file.Close()
    
        // Decode the image format
        _, format, err := image.DecodeConfig(file)
        if err != nil {
            fmt.Println("Error decoding image:", err)
            return
        }
    
        // Check if the image is in PNG format
        if format == "png" {
            fmt.Println("The image is in PNG format.")
        } else {
            fmt.Println("The image is not in PNG format.")
        }
    }


<br><br>
convert a PNG image to JPEG format using the image and image/jpeg packages from the Golang standard library.
    

    package main
    
    import (
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "os"
    )
    
    func main() {
    // Open the PNG file
    pngFile, err := os.Open("input.png")
    if err != nil {
    fmt.Println("Error opening PNG file:", err)
    return
    }
    defer pngFile.Close()
    
        // Decode the PNG image
        pngImage, _, err := image.Decode(pngFile)
        if err != nil {
            fmt.Println("Error decoding PNG image:", err)
            return
        }
    
        // Create a new JPEG file
        jpegFile, err := os.Create("output.jpg")
        if err != nil {
            fmt.Println("Error creating JPEG file:", err)
            return
        }
        defer jpegFile.Close()
    
        // Encode the PNG image as JPEG and write to the new file
        err = jpeg.Encode(jpegFile, pngImage, &jpeg.Options{Quality: 100})
        if err != nil {
            fmt.Println("Error encoding PNG to JPEG:", err)
            return
        }
    
        fmt.Println("Conversion successful. JPEG file created.")
    }

<br><br>
To upload the compressed image to an AWS S3 bucket, use the aws-sdk-go library. 

    package main

    import (
    "bytes"
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/credentials"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/s3"
        "github.com/labstack/echo/v4"
    )
    
    const (
    awsAccessKey = "YOUR_AWS_ACCESS_KEY"
    awsSecretKey = "YOUR_AWS_SECRET_KEY"
    awsRegion    = "YOUR_AWS_REGION"
    s3Bucket     = "YOUR_S3_BUCKET"
    )
    
    func main() {
    e := echo.New()
    
        // Define routes
        e.POST("/compress/jpeg", compressJPEG)
        e.POST("/compress/png", compressPNG)
    
        // Start the server
        e.Logger.Fatal(e.Start(":8080"))
    }
    
    func compressJPEG(c echo.Context) error {
    return compressImage(c, "jpeg")
    }
    
    func compressPNG(c echo.Context) error {
    return compressImage(c, "png")
    }
    
    func compressImage(c echo.Context, format string) error {
    // Get the image file from the request
    file, err := c.FormFile("image")
    if err != nil {
    return c.String(http.StatusBadRequest, "Failed to get image from the request.")
    }
    
        // Open the file
        src, err := file.Open()
        if err != nil {
            return c.String(http.StatusInternalServerError, "Failed to open the image file.")
        }
        defer src.Close()
    
        // Decode the image
        img, _, err := image.Decode(src)
        if err != nil {
            return c.String(http.StatusInternalServerError, "Failed to decode the image.")
        }
    
        // Compress the image
        var compressedImage []byte
        switch format {
        case "jpeg":
            compressedImage, err = encodeJPEG(img)
        case "png":
            compressedImage, err = encodePNG(img)
        default:
            return c.String(http.StatusBadRequest, "Invalid image format. Supported formats: jpeg, png.")
        }
    
        if err != nil {
            return c.String(http.StatusInternalServerError, "Failed to encode the compressed image.")
        }
    
        // Upload the compressed image to AWS S3
        uploadedURL, err := uploadToS3(file.Filename, compressedImage)
        if err != nil {
            return c.String(http.StatusInternalServerError, "Failed to upload the compressed image to S3.")
        }
    
        // Return the S3 URL
        return c.String(http.StatusOK, fmt.Sprintf("Compressed image uploaded to: %s", uploadedURL))
    }
    
    func encodeJPEG(img image.Image) ([]byte, error) {
    var buf bytes.Buffer
    w := io.Writer(&buf)
    err := jpeg.Encode(w, img, &jpeg.Options{Quality: 50})
    return buf.Bytes(), err
    }
    
    func encodePNG(img image.Image) ([]byte, error) {
    var buf bytes.Buffer
    w := io.Writer(&buf)
    err := png.Encode(w, img)
    return buf.Bytes(), err
    }
    
    func uploadToS3(filename string, data []byte) (string, error) {
    // Create a new AWS session
    sess, err := session.NewSession(&aws.Config{
    Region:      aws.String(awsRegion),
    Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
    })
    if err != nil {
    return "", err
    }
    
        // Create an S3 service client
        svc := s3.New(sess)
    
        // Prepare the S3 object key
        objectKey := fmt.Sprintf("uploads/%s", filename)
    
        // Upload the file to S3
        _, err = svc.PutObject(&s3.PutObjectInput{
            Bucket:      aws.String(s3Bucket),
            Key:         aws.String(objectKey),
            ContentType: aws.String(http.DetectContentType(data)),
            Body:        bytes.NewReader(data),
        })
        if err != nil {
            return "", err
        }
    
        // Generate the S3 URL for the uploaded file
        uploadedURL := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", s3Bucket, awsRegion, objectKey)
    
        return uploadedURL, nil
    }



<br><br>
Comparison between JPEG and PNG for compression

|                       | JPEG                                      | PNG                                       |
|-----------------------|-------------------------------------------|-------------------------------------------|
| **Advantages**        |                                           |                                           |
| Compression           | High compression, smaller file size         | Lossless compression, good for simple images and graphics with large areas of uniform color|
| Color Depth           | Supports millions of colors (24-bit)       | Supports millions of colors (24-bit)       |
| Transparency          | Does not support transparency              | Supports alpha channel for transparency   |
| **Disadvantages**     |                                           |                                           |
| Lossy Compression     | Lossy compression, may result in artifacts  | Lossless compression, larger file sizes    |
| Artifacts             | Compression artifacts may be visible       | No artifacts, but larger file sizes        |
| Transparency          | Does not support transparency              | Supports transparency, but may result in larger file sizes|
| Suitable for          | Photographs and images with gradients      | Images with sharp edges, logos, icons      |
