package app

import "github.com/deadblue/fm.rss/internal/db"

type _ItemSlice []*db.Item

func (s _ItemSlice) Len() int {
	return len(s)
}

func (s _ItemSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s _ItemSlice) Less(i, j int) bool {
	if s[i] == nil {
		return false
	}
	if s[j] == nil {
		return true
	}
	return s[i].CreateTime.After(s[j].CreateTime)
}
