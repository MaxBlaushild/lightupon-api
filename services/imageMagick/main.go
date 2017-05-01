package imageMagick

import (
	"gopkg.in/gographics/imagick.v2/imagick"
	"fmt"
)


func CropPin(imageBinary []byte) (transformedBinary []byte) {
	imagick.Initialize()
	defer imagick.Terminate()

	var err error

	mw := imagick.NewMagickWand()

	err = mw.ReadImageBlob(imageBinary)

	if err != nil {
		fmt.Println(err)
	}

	err = mw.ResizeImage(120, 120, imagick.FILTER_LANCZOS, 1)

	if err != nil {
		fmt.Println(err)
	}

	err = mw.SetImageCompressionQuality(100)

	if err != nil {
		fmt.Println(err)
	}

	transformedBinary = mw.GetImageBlob()
	return
}