package facebook

import (
	fb "github.com/huandu/facebook"
	"fmt"
)

type User struct {
	ID string
	AccessToken string
}

type Post struct {
	Message string
	PictureUrl string
	Link string
}

func CreatePost(user User, post Post) (err error) {
	fmt.Println("creating post")
	_, err = fb.Post("/" + user.ID + "/feed", fb.Params{
		"message": post.Message,
		"access_token": user.AccessToken,
		"picture": post.PictureUrl,
		"link": post.Link,
	})
	fmt.Println(err)
	return
}