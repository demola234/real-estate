package utils

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadAvatar(ctx context.Context, image multipart.File) (string, error) {

	cloudinary_url := GoDotEnvVariable("CLOUDINARY_API_KEY")
	cld, _ := cloudinary.NewFromURL(cloudinary_url)
	// Get the preferred name of the file if its not supplied
	fileName := "profileImage"

	result, error := cld.Upload.Upload(ctx, image, uploader.UploadParams{
		PublicID: fileName,
		// Split the tags by comma
		Tags: strings.Split(",", "profile"),
	})

	return result.SecureURL, error

}
