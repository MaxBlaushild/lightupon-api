package models

import (
      _ "github.com/lib/pq"
      _ "github.com/jinzhu/gorm/dialects/postgres"
      "github.com/jinzhu/gorm"
      "log"
      "os"
      "fmt"
)

var (
  DB *gorm.DB
)

func getDatabaseString() (dbString string) {
  dbString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
    os.Getenv("LIGHTUPON_DB_HOST"),
    os.Getenv("LIGHTUPON_DB_PORT"),
    os.Getenv("LIGHTUPON_DB_USERNAME"),
    os.Getenv("LIGHTUPON_DB_NAME"),
    os.Getenv("LIGHTUPON_DB_PASSWORD"))

  return
}

func Connect() {
  var err error

  DB, err = gorm.Open("postgres", getDatabaseString())
  if err != nil {
      log.Fatalln(err)
  }

  DB.LogMode(false)
  DB.AutoMigrate(&User{}, 
                 &Trip{}, 
                 &Party{}, 
                 &Scene{}, 
                 &Card{}, 
                 &Partyuser{}, 
                 &PartyInvite{}, 
                 &Location{}, 
                 &Follow{}, 
                 &TripLike{},
                 &Comment{},
                 &SceneLike{},
                 &Device{},
                 &DiscoveredScene{},
                 &Flag{},
                 &BlacklistUser{},
                 &Vote{},
                 &DiscoveredPost{},
                 &Post{},
                 &Pin{})
  
  DB.Model(&Partyuser{}).AddUniqueIndex("idx_partyuser", "party_id", "user_id")
}