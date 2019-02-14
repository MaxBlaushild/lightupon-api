package models

import(
      "github.com/jinzhu/gorm"
      )

type Scene struct {
  gorm.Model
  Name string
  Latitude float64
  Longitude float64
  TripID uint `gorm:"index"`
  Trip Trip
  BackgroundUrl string `gorm:"not null"`
  SceneOrder uint `gorm:"not null"`
  Cards []Card
  GooglePlaceID string
  Route string
  ShareOnFacebook bool
  User User
  UserID uint
  FormattedAddress string
  Locality string
  Neighborhood string
  PostalCode string
  Country string
  AdministrativeLevelTwo string
  AdministrativeLevelOne string
  StreetNumber string
  SoundKey string
  SoundResource string
  PinUrl string
  SelectedPinUrl string
  ConstellationPoint ConstellationPoint
  Liked bool `sql:"-"`
  PercentDiscovered float64 `sql:"-"`
  RawScore float64 `sql:"-"`
  TimeVoteScore float64 `sql:"-"`
  SpatialScore float64 `sql:"-"`
}