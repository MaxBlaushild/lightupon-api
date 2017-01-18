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

// If this is run with testMode=true, the global DB variable will point at the test database
func Connect(testMode bool) {

  var dbString string
  var err error

  if !(testMode) {
    dbString = os.Getenv("DATABASE_URL")  // Not sure why we have 2 ways of getting the DB url. Should probably fix that...
    if len(dbString) == 0 {
      dbString = "user=" + os.Getenv("DB_USERNAME") + " dbname=" + os.Getenv("DB_NAME") + " sslmode=disable"
    }
  } else {
    dbString = "user=" + os.Getenv("DB_USERNAME") + " dbname=" + os.Getenv("DB_TEST_NAME") + " sslmode=disable"
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
                 &SceneLike{})
  
  DB.Model(&Partyuser{}).AddUniqueIndex("idx_partyuser", "party_id", "user_id")
}