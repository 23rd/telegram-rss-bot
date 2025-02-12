package chans

import (
	"fmt"
	"github.com/0x111/telegram-rss-bot/conf"
	"github.com/0x111/telegram-rss-bot/feeds"
	"github.com/0x111/telegram-rss-bot/replies"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

// Get Feed Updates from the feeds
func FeedUpdates() {
	feedUpdates := feeds.GetFeedUpdatesChan()

	for feedUpdate := range feedUpdates {
		log.WithFields(log.Fields{"feedData": feedUpdate}).Debug("Requesting feed data")
		log.WithFields(log.Fields{"feedID": feedUpdate.ID, "feedUrl": feedUpdate.Url}).Info("Updating feeds")
		feeds.GetFeed(feedUpdate.Url, feedUpdate.ID)
	}
}

// Post Feed data to the channel
func FeedPosts(Bot *tgbotapi.BotAPI) {
	feedPosts := feeds.PostFeedUpdatesChan()
	feedJoin := conf.GetConfig().GetString("feed_join")
	if feedJoin == "" {
		feedJoin = " - "
	}

	for feedPost := range feedPosts {
		msg := fmt.Sprintf(`
	%s%s%s
	`, feedPost.Title, feedJoin, feedPost.Link)
		log.WithFields(log.Fields{"feedPost": feedPost, "chatID": feedPost.ChatID}).Debug("Posting feed update to the Telegram API")
		err := replies.SimpleMessage(Bot, feedPost.ChatID, 0, msg)
		if err == nil {
			log.WithFields(log.Fields{"feedPost": feedPost, "chatID": feedPost.ChatID}).Debug("Setting the Feed Data entry to published!")
			_, err := feeds.UpdateFeedDataPublished(&feedPost)
			if err != nil {
				log.WithFields(log.Fields{"error": err, "feedPost": feedPost, "chatID": feedPost.ChatID}).Error("There was an error while updating the Feed Data entry!")
			}
		} else {
			log.WithFields(log.Fields{"error": err, "feedPost": feedPost, "chatID": feedPost.ChatID}).Error("There was an error while posting the update to the feed!")
		}
	}
}
