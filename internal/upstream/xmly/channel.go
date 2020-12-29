package xmly

import (
	"fmt"
	"github.com/deadblue/fm.rss/internal/db"
	"time"
)

type _AlbumData struct {
	Album struct {
		Title          string `json:"title"`
		Intro          string `json:"intro"`
		Cover          string `json:"detailCoverPath"`
		Author         string `json:"nickname"`
		CreateTime     int64  `json:"createdAt"`
		UpdateTime     int64  `json:"updatedAt"`
		TrackCount     int    `json:"tracks"`
		SerializeState int    `json:"serializeStatus"`
	} `json:"album"`
}

const (
	urlAlbumData = "http://mobile.ximalaya.com/mobile-album/album/page/ts-%d?device=android" +
		"&isQueryInvitationBrand=true&albumId=%s&isAsc=true&isVideoAsc=true" +
		"&pageId=1&pageSize=20&pre_page=0&source=0&supportWebp=true"
)

func (f *fetcherImpl) FetchChannel(channelId string, channel *db.Channel) (hasUpdate bool, err error) {
	// Fetch album data from upstream
	url := fmt.Sprintf(urlAlbumData, time.Now().Unix()*1000, channelId)
	data := &_AlbumData{}
	err = f.getJsonV1(url, data)
	if err != nil {
		return
	}
	// Fill all fields for new document
	updateTime, now := time.Unix(data.Album.UpdateTime/1000, 0), time.Now()
	if channel.ChannelId == "" {
		channel.Source = _Name
		channel.ChannelId = channelId
		channel.Link = makeChannelLink(channelId)
		channel.Title = data.Album.Title
		channel.Intro = data.Album.Intro
		channel.Cover = data.Album.Cover
		channel.Author = data.Album.Author
		channel.ItemCount = data.Album.TrackCount
		channel.IsComplete = data.Album.SerializeState == 2
		channel.CreateTime = time.Unix(data.Album.CreateTime/1000, 0)
		channel.UpdateTime = updateTime
		channel.FetchTime = now
		hasUpdate = true
	} else if channel.UpdateTime.Before(updateTime) {
		channel.UpdateTime = updateTime
		channel.FetchTime = now
		hasUpdate = true
	}
	return
}

func makeChannelLink(channelId string) string {
	return fmt.Sprintf("https://www.ximalaya.com/album/%s/", channelId)
}
