package upstream

import "github.com/deadblue/fm.rss/internal/db"

type Fetcher interface {

	// Name returns the upstream source name.
	Name() string

	// Aliases returns all upstream alias names.
	Aliases() []string

	// FetchChannel fetches channel data from upstream, and fills the document.
	FetchChannel(channelId string, channel *db.Channel) (hasUpdate bool, err error)

	// FetchChannelItems fetches all items' IDs under the channel.
	FetchChannelItems(channelId string) (itemIds []string, err error)

	// FetchItems fetches items data and fills the documents.
	FetchItems(channelId string, items []*db.Item) (err error)
}
