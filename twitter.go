package htmlhouse

import (
	"database/sql"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
)

func tweet(app *app, houseID, title string) {
	// Check for blacklisted titles
	if title == "HTMLhouse" {
		return
	}

	// Check if this has already been tweeted
	var tweetID int64
	err := app.db.QueryRow("SELECT tweet_id FROM tweetedhouses WHERE house_id = ?", houseID).Scan(&tweetID)
	switch {
	case err != nil && err != sql.ErrNoRows:
		fmt.Printf("Error selecting from tweetedhouses: %v", err)
		return
	}
	if tweetID != 0 {
		return
	}

	// Post to Twitter
	text := fmt.Sprintf("\"%s\" on #HTMLhouse - %s/%s.html #html #web #website", title, app.cfg.HostName, houseID)

	anaconda.SetConsumerKey(app.cfg.TwitterConsumerKey)
	anaconda.SetConsumerSecret(app.cfg.TwitterConsumerSecret)
	api := anaconda.NewTwitterApi(app.cfg.TwitterToken, app.cfg.TwitterTokenSecret)

	t, err := api.PostTweet(text, nil)
	if err != nil {
		fmt.Printf("Error posting tweet: %v", err)
	}

	// Mark it as "tweeted"
	_, err = app.db.Exec("INSERT INTO tweetedhouses (house_id, tweet_id) VALUES (?, ?)", houseID, t.Id)
	if err != nil {
		fmt.Printf("Error noting house tweet status: %v", err)
		return
	}
}
