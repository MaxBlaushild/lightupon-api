package models

import (
	     "github.com/jinzhu/gorm"
      "net/http"
      "io/ioutil"
      "lightupon-api/services/aws"
      "lightupon-api/services/imageMagick"
      "fmt"
)

type Pin struct {
	gorm.Model
  Url string
  OwnerID uint
  OwnerType string
}

func NewPin(url string, id uint, ownerType string) (pin Pin, err error) {
	binary, err := DownloadImage(url)
	pinBinary := imageMagick.CropPin(binary)

  asset := aws.Asset{
    Type: "images", 
    Name: assetName(id, ownerType), 
    Extension: ".png",
    Binary: pinBinary,
  }

  getUrl, err := aws.UploadAsset(asset); if err != nil {
  	return
  }

  pin = Pin{
  	OwnerID: id,
  	OwnerType: ownerType,
  	Url: getUrl,
  }

  err = DB.Create(&pin).Error

  return
}

func assetName(id uint, ownerType string) string {
  return "/" + ownerType  + "/" + fmt.Sprintf("%v", id) + "/pin"
}

func DownloadImage(url string) (imageBinary []byte, err error) {
  resp, err := http.Get(url)

  defer resp.Body.Close()

  imageBinary, err = ioutil.ReadAll(resp.Body); if err != nil {
    fmt.Println("ioutil.ReadAll -> %v", err)
  }

  return
}
