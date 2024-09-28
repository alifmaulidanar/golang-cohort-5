package config

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

// InitializeCloudinary initializes Cloudinary client
func InitializeCloudinary() (*cloudinary.Cloudinary, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get Cloudinary credentials from environment variables
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return cld, nil
}

// UploadImage uploads an image to Cloudinary
func UploadImage(cld *cloudinary.Cloudinary, filePath string, folder string) (*uploader.UploadResult, error) {
	ctx := context.Background()

	uploadParams := uploader.UploadParams{
		Folder: folder,
	}

	result, err := cld.Upload.Upload(ctx, filePath, uploadParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}
