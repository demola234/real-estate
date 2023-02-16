package utils

import (
	"context"
	"testing"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/stretchr/testify/require"
)

func TestUploadAvatar(t *testing.T) {
	cld, _ := cloudinary.NewFromURL("cloudinary://211576879732455:W6p_HMMIrDZkEfheHRUHIkSTdOo@dcnuiaskr")
	image := "C:\\Users\\User\\Desktop\\profile.jpg"
	// Get the preferred name of the file if its not supplied
	fileName := RandomFileName()

	_, error := cld.Upload.Upload(context.Background(), image, uploader.UploadParams{
		PublicID: fileName,
	})
	require.NoError(t, error)
}
