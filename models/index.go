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

func getDatabaseString(productionMode bool) (dbString string) {
  if productionMode {
     dbString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
      os.Getenv("LIGHTUPON_DB_HOST"),
      os.Getenv("LIGHTUPON_DB_PORT"),
      os.Getenv("LIGHTUPON_DB_USERNAME"),
      os.Getenv("LIGHTUPON_DB_NAME"),
      os.Getenv("LIGHTUPON_DB_PASSWORD"))
  } else {  
    dbString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
      os.Getenv("LIGHTUPON_DB_HOST"),
      os.Getenv("LIGHTUPON_DB_PORT"),
      os.Getenv("LIGHTUPON_DB_USERNAME"),
      os.Getenv("LIGHTUPON_TEST_DB_NAME"),
      os.Getenv("LIGHTUPON_DB_PASSWORD"))
  }

  return
}

func Connect(productionMode bool) {
  var err error

  DB, err = gorm.Open("postgres", getDatabaseString(productionMode))
  if err != nil {
      log.Fatalln(err)
  }

  DB.LogMode(false)
  DB.AutoMigrate(&User{},
                 &Location{},
                 &Device{},
                 &Flag{},
                 &BlacklistUser{},
                 &DiscoveredPost{},
                 &Post{},
                 &Pin{})

  DatabaseUpdateTemporary() // This will update fields that need to be updated in order for things to work
}