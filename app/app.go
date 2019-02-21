package app

import(
  "lightupon-api/models"
)

func GetNearbyPosts(userID uint, lat string, lon string, radius string, numPosts int, databaseAccessor models.DatabaseAccessor) (posts []models.Post, err error) {
  nearbyTipPosts, _ := getNearbyTipPosts(userID, lat, lon, radius, databaseAccessor)
  posts = append(posts, nearbyTipPosts...)

  nearbyCompletedPosts, _ := databaseAccessor.GetNearbyCompletedPosts(userID, lat, lon, radius)
  posts = append(posts, nearbyCompletedPosts...)

  nearbyUncompletedFirstPosts, _ := databaseAccessor.GetNearbyUncompletedFirstPosts(userID, lat, lon, radius)
  posts = append(posts, nearbyUncompletedFirstPosts...)

  return
}

func getNearbyTipPosts(userID uint, lat string, lon string, radius string, databaseAccessor models.DatabaseAccessor) (tipPosts []models.Post, err error) {
  // This will get the quest_order and quest_id for the maximum completed post in each quest for the user.
  results, _ := databaseAccessor.GetQuestOrderForLastCompletedPostInEachQuest(userID)

  var post models.Post

  for _, result := range results {
    // Let's try to get the very next post in the quest (only if it's nearby). If we can't find one, then the user has completed the entire quest or the tip post is not nearby.
    post, _ = databaseAccessor.FindNearbyPostInQuestWithParticularQuestOrder(lat, lon, radius, result.QuestID, result.MaxQuestOrder + 1)
    if post.ID != 0 {
      tipPosts = append(tipPosts, post)
    }
  }

  return
}