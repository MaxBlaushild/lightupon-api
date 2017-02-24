package imageMagick

import (
	"github.com/quirkey/magick"
	"log"
)


func CropPin(imageBinary []byte) (transformedBinary []byte) {
	image, err := magick.NewFromBlob(imageBinary, "png")
	err = image.Resize("40x40"); if err != nil {
		log.Print("Problem with transforming")
	}

	err = image.Shadow("#F00", 255, 5, 2, 2); if err != nil {
		log.Print("Problem with transforming")
	}

	err = image.FillBackgroundColor("#00F"); if err != nil {
		log.Print("Problem setting background")
	}

	transformedBinary, err = image.ToBlob("png")
	return
}