package routes

import (
				"lightupon-api/aws"
				"net/http"
				"github.com/gorilla/mux"
				)

func UploadAssetUrlHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  assetType := vars["assetType"]
  assetName := vars["assetName"]
  urlStr, err := aws.PutAsset(assetType, assetName)

  if (err == nil) {
  	respondWithAccepted(w, urlStr)
  } else {
  	respondWithBadRequest(w, "The type provided in the uri is invalid.")
  }
}

func GetAssetUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  assetType := vars["assetType"]
  assetName := vars["assetName"]
  urlStr, err := aws.GetAsset(assetType, assetName)

  if (err == nil) {
  	respondWithAccepted(w, urlStr)
  } else {
  	respondWithBadRequest(w, "The type provided in the uri is invalid.")
  }
}
