package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/gorilla/mux"
       "github.com/jinzhu/gorm"
       "text/template"
       )

func AdminSceneDetailsHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  scene_id := vars["id"]

  t := template.New("fieldname example")
  t, _ = t.Parse(scene_detail_template)
  scene := models.Scene{}
  models.DB.Preload("Cards", func(DB *gorm.DB) *gorm.DB {
    return DB.Order("Cards.card_order ASC")
  }).Where("id = $1", scene_id).Find(&scene)
  t.Execute(w, scene)
}

func AdminTripDetailsHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  trip_id := vars["id"]

  t := template.New("fieldname example")
  t, _ = t.Parse(trip_detail_template)
  trip := models.Trip{}
  models.DB.Preload("Scenes").Where("id = $1", trip_id).Find(&trip)
  t.Execute(w, trip)
}

func AdminGetTripsForUserHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]

  t := template.New("fieldname example")
  t, _ = t.Parse(trips_list_template)
  trips := []models.Trip{}
  models.DB.Where("owner = $1", id).Find(&trips)
  t.Execute(w, trips)
}

func AdminGetAllTripsHandler(w http.ResponseWriter, r *http.Request) {
  t := template.New("fieldname example")
  t, _ = t.Parse(trips_list_template)
  trips := []models.Trip{}
  models.DB.Find(&trips)
  t.Execute(w, trips)
}