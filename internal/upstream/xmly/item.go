package xmly

import (
	"fmt"
	"github.com/deadblue/fm.rss/internal/db"
	"log"
	"strconv"
	"strings"
	"time"
)

type _TrackData struct {
	TrackId    int    `json:"trackId"`
	Title      string `json:"title"`
	Duration   int    `json:"duration"`
	CreateTime int64  `json:"createdAt"`
}

type _AlbumTracksData struct {
	List       []*_TrackData `json:"list"`
	PageSize   int           `json:"pageSize"`
	PageId     int           `json:"pageId"`
	MaxPageId  int           `json:"maxPageId"`
	TotalCount int           `json:"totalCount"`
}

type _TrackDetailData struct {
	Title    string `json:"title"`
	Intro    string `json:"intro"`
	Duration int    `json:"duration"`

	LqUrl   string `json:"playUrl32"`
	LqSize  int64  `json:"playUrl32Size"`
	Hq1Url  string `json:"playPathAacv164"`
	Hq1Size int64  `json:"playPathAacv164Size"`
	Hq2Url  string `json:"playPathAacv224"`
	Hq2Size int64  `json:"playPathAacv224Size"`

	CreateTime int64 `json:"createdAt"`
}

const (
	urlAlbumTracks = "http://mobile.ximalaya.com/mobile/v1/album/track/ts-%d?device=android" +
		"&albumId=%s&isAsc=true&pageSize=20&pageId=%d"
	urlTrackDetail = "http://mobile.ximalaya.com/mobile/track/v2/baseInfo/%d?device=android&trackId=%s"
)

func (f *fetcherImpl) FetchChannelItems(channelId string) (itemIds []string, err error) {
	itemIds = make([]string, 0)
	for pageId := 1; ; pageId++ {
		// Fetch track data
		url := fmt.Sprintf(urlAlbumTracks, time.Now().Unix()*1000, channelId, pageId)
		data := &_AlbumTracksData{}
		if err = f.getJsonV1(url, data); err != nil {
			log.Printf("Get tracks at page %d error: %s", pageId, err)
			break
		}
		for _, track := range data.List {
			itemIds = append(itemIds, strconv.Itoa(track.TrackId))
		}
		if pageId == data.MaxPageId {
			break
		}
	}
	return
}

func (f *fetcherImpl) FetchItems(channelId string, items []*db.Item) (err error) {
	for _, item := range items {
		url := fmt.Sprintf(urlTrackDetail, time.Now().Unix()*1000, item.ItemId)
		data := &_TrackDetailData{}
		if err = f.getJsonV2(url, "trackInfo", data); err != nil {
			log.Printf("Fetch item [%s/%s] detail error: %s", _Name, item.ItemId, err)
			return
		}
		item.Link = makeItemLink(channelId, item.ItemId)
		item.Title = data.Title
		item.Intro = data.Intro
		item.Duration = data.Duration
		item.CreateTime = time.Unix(data.CreateTime/1000, 0)
		item.Url = data.Hq2Url
		item.Size = data.Hq2Size
		if strings.HasSuffix(item.Url, ".mp3") {
			item.MimeType = "audio/mpeg"
		} else if strings.HasSuffix(item.Url, ".aac") {
			item.MimeType = "audio/aac"
		} else if strings.HasSuffix(item.Url, ".m4a") {
			item.MimeType = "audio/x-m4a"
		}
		item.FetchTime = time.Now()
	}
	return
}

func makeItemLink(channelId, itemId string) string {
	return fmt.Sprintf("https://www.ximalaya.com/track/%s/%s/", channelId, itemId)
}
