package imageMagick

import (
	"github.com/quirkey/magick"
	"log"
)


func CropPin(imageBinary []byte, size string) (transformedBinary []byte) {
	image, err := magick.NewFromBlob(imageBinary, "png")
	defer image.Destroy()

	err = image.Resize(size); if err != nil {
		log.Print("Problem with transforming")
	}

	transformedBinary, err = image.ToBlob("png")
	return
}