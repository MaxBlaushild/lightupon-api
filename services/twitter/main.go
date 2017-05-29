package twitter

import(
			"github.com/ChimeraCoder/anaconda"
			"os"
			"net/url"
			"fmt"
			"strconv"
			"encoding/base64"
			// "github.com/kr/pretty"
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
	MediaID int64
}

func newClient(user User) *anaconda.TwitterApi {
	return anaconda.NewTwitterApi(user.AccessToken, user.AccessTokenSecret)
}

func PostStatus(user User, status Status) (err error) {
	client := newClient(user)
	values := url.Values{}
	latString := fmt.Sprintf("%.6f", status.Lat)
  longString := fmt.Sprintf("%.6f", status.Long)
  mediaIDString := strconv.FormatInt(status.MediaID, 10)
  mediaIDString = "[" + mediaIDString
  mediaIDString = mediaIDString + "]"
	values.Set("lat", latString)
	values.Set("long", longString)
	values.Set("media_ids", mediaIDString)
	_, err = client.PostTweet(status.Status, values)
	return
}

func PostMedia(user User, mediaBinary []byte) (media anaconda.Media, err error) {
	client := newClient(user)
	imgBase64Str := base64.StdEncoding.EncodeToString(mediaBinary)
	media, err = client.UploadMedia(imgBase64Str)
	return
}