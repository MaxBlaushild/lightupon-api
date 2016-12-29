package routes

import (
				"lightupon-api/services/aws"
        "lightupon-api/models"
				"net/http"
				"github.com/gorilla/mux"
        "encoding/json"
				)

func UploadAssetUrlHandler(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  asset := models.Asset{}
  err := decoder.Decode(&asset)
  urlStr, err := aws.PutAsset(asset.Type, asset.Name)

  if (err == nil) {
  	respondWithAccepted(w, urlStr)
  } else {
  	respondWithBadRequest(w, "The type provided in the uri is invalid.")
  }
}

func GetAssetUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  assetType := vars["type"]
  assetName := vars["name"]
  urlStr, err := aws.GetAsset(assetType, assetName)

  if (err == nil) {
  	respondWithAccepted(w, urlStr)
  } else {
  	respondWithBadRequest(w, "The type provided in the uri is invalid.")
  }
}