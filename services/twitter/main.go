package twitter

import(
			"github.com/ChimeraCoder/anaconda"
			"os"
			"net/url"
			"fmt"
)

func Init() {
	anaconda.SetConsumerKey(os.Getenv("LIGHTUPON_TWITTER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("LIGHTUPON_TWITTER_SECRET"))
}

type User struct {
	AccessToken string
	AccessTokenSecret string
}

type Status struct {
	Lat float64
	Long float64
	Status string
}

func newClient(user User) *anaconda.TwitterApi {
	return anaconda.NewTwitterApi(user.AccessToken, user.AccessTokenSecret)
}

func PostStatus(user User, status Status) (err error) {
	client := newClient(user)
	values := url.Values{}
	latString := fmt.Sprintf("%.6f", status.Lat)
  longString := fmt.Sprintf("%.6f", status.Long)
	values.Set("lat", latString)
	values.Set("long", longString)
	_, err = client.PostTweet(status.Status, values)
	return
}