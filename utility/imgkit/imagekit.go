package imgkit

import (
	"context"
	"encoding/json"
	"fmt"
	"libsysfo-server/database"
	"os"
	"time"

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
	t := time.Now().UnixMilli()
	filename := fmt.Sprintf("%d-%s", t, data.FileName)
	ur := imagekit.UploadRequest{
		File:              data.File,
		FileName:          filename,
		UseUniqueFileName: false,
		Tags:              []string{"testing", "test"},
		Folder:            data.Folder,
		IsPrivateFile:     false,
		CustomCoordinates: "",
		ResponseFields:    nil,
	}

	ctx := context.Background()

	upr, err = ik.Upload.ServerUpload(ctx, &ur)
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := json.Marshal(upr)
	if err != nil {
		fmt.Println(err)
		return
	}
	database.DB.Save(&database.ThirdPartyJobs{
		Job:          "upload Image",
		Destination:  "ImageKit",
		ResponseBody: string(b),
		Status:       200,
	})
	return
}

type ImgInformation struct {
	File     interface{}
	FileName string
	Folder   string
}
