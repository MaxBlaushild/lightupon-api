package twitter

import(
			"github.com/ChimeraCoder/anaconda"
			"os"
			"net/url"
			"fmt"
			"encoding/base64"
			"github.com/kr/pretty"
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
	MediaID string
}

func newClient(user User) *anaconda.TwitterApi {
	return anaconda.NewTwitterApi(user.AccessToken, user.AccessTokenSecret)
}

func PostStatus(user User, status Status) (err error) {
	client := newClient(user)
	values := url.Values{}
	latString := fmt.Sprintf("%.6f", status.Lat)
  longString := fmt.Sprintf("%.6f", status.Long)
  mediaIDString := status.MediaID
  mediaIDString = "[" + mediaIDString
  mediaIDString = mediaIDString + "]"
  pretty.Println(mediaIDString)
	values.Set("lat", latString)
	values.Set("long", longString)
	values.Set("media_ids", mediaIDString)
	res, err := client.PostTweet(status.Status, values)
	pretty.Println(res)
	pretty.Println(err)
	return
}

func PostMedia(user User, mediaBinary []byte) (media anaconda.Media, err error) {
	client := newClient(user)
	imgBase64Str := base64.StdEncoding.EncodeToString(mediaBinary)
	media, err = client.UploadMedia(imgBase64Str)
	pretty.Println(media)
	return
}