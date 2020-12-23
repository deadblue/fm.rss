package app

import (
	"errors"
	"fmt"
	"github.com/deadblue/fm.rss/internal/db"
	"github.com/deadblue/fm.rss/internal/rss"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
)

var (
	_RegexFeedPath = regexp.MustCompile("^/(\\w+)/(\\w+)\\.xml$")

	errUnsupportedSource = errors.New("unsupported source")

	_ErrMalformedUrl = []byte("Malformed URL")
)

func (a *appImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matches := _RegexFeedPath.FindStringSubmatch(r.URL.Path)
	if matches == nil || len(matches) == 0 {
		// Show error
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(_ErrMalformedUrl)
	} else {
		feed, err := a.requestFeed(matches[1], matches[2])
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		} else {
			content := []byte(feed.String())
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
			w.Header().Set("Content-Length", strconv.Itoa(len(content)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(content)
		}
	}
}

func (a *appImpl) requestFeed(source, channelId string) (feed *rss.Feed, err error) {
	f, ok := a.fs[source]
	if !ok {
		return nil, errUnsupportedSource
	}
	source, feed = f.Name(), rss.New()
	log.Printf("Request feed [%s/%s] ...", source, channelId)
	// Retrieve channel
	channel, isNew, hasUpdate := &db.Channel{}, false, false
	if err = a.dc.ChannelRetrieve(source, channelId, channel); err != nil {
		if err == db.ErrNotFound {
			err, isNew = nil, true
		} else {
			return
		}
	}
	if !channel.IsComplete {
		log.Printf("Fetching channel [%s/%s] from upstream ...", source, channelId)
		if hasUpdate, err = f.FetchChannel(channelId, channel); err != nil {
			return
		}
		if isNew {
			log.Printf("Saving new channel [%s/%s] to database ...", source, channelId)
			err = a.dc.ChannelInsert(channel)
		} else if hasUpdate {
			log.Printf("Updating channel [%s/%s] in database ...", source, channelId)
			err = a.dc.ChannelUpdate(channel)
		}
		if err != nil {
			return
		}
	}
	// Fill channel info feed
	feed.WithInfo(channel.Link, channel.Title, channel.Intro).
		WithAuthor(channel.Author).
		WithImage(channel.Cover).
		WithPubDate(channel.UpdateTime).
		WithComplete(channel.IsComplete)
	// Retrieve items
	items, err := a.dc.ItemRetrieve(source, channelId)
	if err != nil {
		return
	}
	log.Printf("Cached [%s/%s] items: %d", source, channelId, len(items))
	// Return when no update
	if hasUpdate {
		// Mark cached items
		cachedItems := make(map[string]bool)
		for _, item := range items {
			cachedItems[item.ItemId] = true
		}
		// Fetch all items Id
		var itemIds []string
		if itemIds, err = f.FetchChannelItems(channelId); err != nil {
			log.Printf("Fetch [%s/%s] item list error: %s", source, channelId, err)
			return
		}
		// Create new items
		newItems := make([]*db.Item, 0)
		for _, itemId := range itemIds {
			if cachedItems[itemId] {
				continue
			} else {
				newItems = append(newItems, &db.Item{
					Source:    source,
					ItemId:    itemId,
					ChannelId: channelId,
					Guid:      fmt.Sprintf("%s-%s-%s", source, channelId, itemId),
				})
			}
		}
		if newCount := len(newItems); newCount > 0 {
			log.Printf("New [%s/%s] items: %d", source, channelId, newCount)
			if err := f.FetchItems(channelId, newItems); err != nil {
				log.Printf("Fetch item detail error: %s", err)
			} else {
				items = append(items, newItems...)
				if err := a.dc.ItemInsert(newItems); err != nil {
					log.Printf("Create item cache error: %s", err)
				}
			}
		}
	}
	// Sort items by time desc
	sort.Sort(_ItemSlice(items))
	for _, item := range items {
		feed.NewItem().
			WithInfo(item.Guid, item.Link, item.Title, item.Intro).
			WithEnclosure(item.Url, item.Size, item.MimeType).
			WithDuration(item.Duration).
			WithPubDate(item.CreateTime)
	}
	return
}
