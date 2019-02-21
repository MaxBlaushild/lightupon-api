package app

import(
       // "net/http"
       "lightupon-api/models"
       // "github.com/davecgh/go-spew/spew"
       // "encoding/json"
       // "github.com/gorilla/mux"
       // "strconv"
       // "fmt"
)

func GetNearbyPosts(userID uint, lat string, lon string, radius string, numPosts int, databaseAccessor models.DatabaseAccessor) (posts []models.Post, err error) {
  tipPosts, _ := getTipPosts(userID, lat, lon, radius, databaseAccessor)
  posts = append(posts, tipPosts...)

  // posts, err = models.GetPostsNearLocation_NEW(lat, lon, userID, radius)

  return
}

func getTipPosts(userID uint, lat string, lon string, radius string, databaseAccessor models.DatabaseAccessor) (tipPosts []models.Post, err error) {
  // This will get the quest_order and quest_id for the maximum completed post in each quest for the user.
  results, _ := databaseAccessor.GetQuestOrderForLastCompletedPostInEachQuest(userID)

  var post models.Post

  for _, result := range results {
    // Let's try to get the very next post in the quest. If we can't find one, then the user has completed the entire quest.
    post, _ = databaseAccessor.FindNearbyPostInQuestWithParticularQuestOrder(lat, lon, radius, result.QuestID, result.MaxQuestOrder + 1)
    if post.ID != 0 {
      tipPosts = append(tipPosts, post)
    }
  }

  return
}