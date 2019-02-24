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
                 &Pin{},
                 &Quest{})

  // setUpTestDataWithoutQuestOrders()
  // setUpTestData()
  DatabaseUpdateTemporary() // This will update fields that need to be updated in order for things to work. Should be removed after it's been run on all machines (dev and prod).
}

func setUpTestDataWithoutQuestOrders() {
  fmt.Println("setUpTestDataWithoutQuestOrders")
  DB.Exec(`DELETE FROM posts;
          DELETE FROM discovered_posts;
          DELETE FROM users;
          DELETE FROM quests;

          INSERT INTO users
          (id)
          VALUES
          (1),
          (2);

          INSERT INTO quests
          (id, description)
          VALUES
          (1, 'This is the dam tour.'),
          (2, 'This is the second dam tour.'),
          (3, 'This is the third dam tour.');

          INSERT INTO posts
          (id, caption, latitude, longitude, image_url)
          VALUES
          (1, 'Caption for scene 1', 42.3439129, -71.0739857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (2, 'Caption for scene 2', 42.3449129, -71.0749857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (3, 'Caption for scene 3', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (4, 'Caption for scene 4', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),

          (5, 'Caption for scene 5', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (6, 'Caption for scene 6', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (7, 'Caption for scene 7', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),

          (8, 'Caption for scene 8', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (9, 'Caption for scene 9', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (10, 'Caption for scene 10', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg');

          INSERT INTO discovered_posts
          (id, user_id, post_id, percent_discovered, completed)
          VALUES
          (1, 1, 1, 1.0, True),
          (2, 1, 2, 1.0, True),
          (3, 1, 3, 1.0, True),
          (4, 1, 4, 1.0, True),

          (5, 1, 5, 1.0, True),
          (6, 1, 6, 0.5, False),

          (9, 2, 1, 0.8, False);`)
}

func setUpTestData() {
  DB.Exec(`DELETE FROM posts;
          DELETE FROM discovered_posts;
          DELETE FROM users;
          DELETE FROM quests;

          INSERT INTO users
          (id)
          VALUES
          (1),
          (2);

          INSERT INTO quests
          (id, description)
          VALUES
          (1, 'This is the dam tour.'),
          (2, 'This is the second dam tour.'),
          (3, 'This is the third dam tour.');

          INSERT INTO posts
          (id, quest_id, quest_order, caption, latitude, longitude, image_url)
          VALUES
          (1, 1, 1, 'Caption for scene 1', 42.3439129, -71.0739857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (2, 1, 2, 'Caption for scene 2', 42.3449129, -71.0749857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (3, 1, 3, 'Caption for scene 3', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (4, 1, 4, 'Caption for scene 4', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),

          (5, 2, 1, 'Caption for scene 5', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (6, 2, 2, 'Caption for scene 6', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (7, 2, 3, 'Caption for scene 7', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),

          (8, 3, 1, 'Caption for scene 8', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (9, 3, 2, 'Caption for scene 9', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg'),
          (10, 3, 3, 'Caption for scene 10', 42.3459129, -71.0759857, 'https://i.ytimg.com/vi/PuCzKf3Hzj0/hqdefault.jpg');

          INSERT INTO discovered_posts
          (id, user_id, post_id, percent_discovered, completed)
          VALUES
          (1, 1, 1, 1.0, True),
          (2, 1, 2, 1.0, True),
          (3, 1, 3, 1.0, True),
          (4, 1, 4, 1.0, True),

          (5, 1, 5, 1.0, True),
          (6, 1, 6, 0.5, False),

          (9, 2, 1, 0.8, False);`)
}