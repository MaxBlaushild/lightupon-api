package routes

import (
				"lightupon-api/services/aws"
				"net/http"
				"github.com/gorilla/mux"
        "encoding/json"
				)

func UploadAssetUrlHandler(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  asset := aws.Asset{}
  err := decoder.Decode(&asset)
  urlStr, err := aws.PutAsset(asset)

  if (err == nil) {
  	respondWithAccepted(w, urlStr)
  } else {
  	respondWithBadRequest(w, "The type provided in the uri is invalid.")
  }
}

func GetAssetUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  asset := aws.Asset{
    Type: vars["type"],
    Name: vars["name"],
  }
  urlStr, err := aws.GetAsset(asset)

  if (err == nil) {
  	respondWithAccepted(w, urlStr)
  } else {
  	respondWithBadRequest(w, "The type provided in the uri is invalid.")
  }
}