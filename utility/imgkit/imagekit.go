package imgkit

import (
	"context"
	"os"

	"github.com/codedius/imagekit-go"
)

func (data ImgInformation) UploadImage() (upr *imagekit.UploadResponse, err error) {
	opts := imagekit.Options{
		PublicKey:  os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		PrivateKey: os.Getenv("IMAGEKIT_PRIVATE_KEY"),
	}

	ik, err := imagekit.NewClient(&opts)
	if err != nil {
		return
	}
	ur := imagekit.UploadRequest{
		File:              data.File,
		FileName:          data.FileName,
		UseUniqueFileName: false,
		Tags:              []string{"testing", "test"},
		Folder:            data.Folder,
		IsPrivateFile:     false,
		CustomCoordinates: "",
		ResponseFields:    nil,
	}

	ctx := context.Background()

	upr, err = ik.Upload.ServerUpload(ctx, &ur)
	return
}

type ImgInformation struct {
	File     interface{}
	FileName string
	Folder   string
}
