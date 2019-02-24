package models

import (
  "testing"
  // "github.com/davecgh/go-spew/spew"
  "fmt"
  )

/*

The following integration test requires the creation of a database called lightupon_test and an environment variable being set called LIGHTUPON_TEST_DB_NAME=lightupon_test

User 1 should get Scenes 1,2,3,4,5,6,8
User 2 should get Scenes 1,5,8

*/

func TestGetNearbyPostsAndTryToDiscoverThem(t *testing.T) {
  Connect(false)
  // Connect(true) // THIS WILL MODIFY THE PRIMARY DATABASE! Only uncomment this if you're Jon and you don't have any data in the primary database that you care about.
  setUpTestData()

  var user User
  DB.Where("id=1").First(&user) // Because of GORM, we're not allowed to set the ID of a user because it's an inherited field. So we have to insert it into the database (done in the test data set up below) and then retrieve it here.

  posts, _ := GetNearbyPostsAndTryToDiscoverThem(user, "42.3459129", "-71.0759857", "5000", 20)
  
  for _, k := range posts {
    fmt.Println(k.ID)
  }

  // TODO: execute for user 2 and also programmatically check that the list of posts is correct
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