package models

import (
      _ "github.com/lib/pq"
      _ "github.com/jinzhu/gorm/dialects/postgres"
      "github.com/jinzhu/gorm"
      "log"
      "os"
)

var (
  DB *gorm.DB
)

func Connect() {
  var dbString string = os.Getenv("DATABASE_URL")
  var err error
  if len(dbString) == 0 {
    dbString = "user=" + os.Getenv("DB_USERNAME") + " dbname=" + os.Getenv("DB_NAME") + " sslmode=disable"
  }

  DB, err = gorm.Open("postgres", dbString)
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
                 &Bookmark{}, 
                 &Like{},
                 &TripLike{},
                 &Comment{},
                 &SceneLike{},
                 &Device{},
                 &ExposedScene{},
                 &NeighborhoodPoint{})
  
  DB.Model(&Partyuser{}).AddUniqueIndex("idx_partyuser", "party_id", "user_id")
}