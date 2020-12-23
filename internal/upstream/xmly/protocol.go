package xmly

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deadblue/gostream/quietly"
	"net/http"
)

const (
	_UserAgent = "ting_6.7.12"
)

type _JsonResponseV1 struct {
	Ret  int             `json:"ret"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (f *fetcherImpl) getJsonV1(url string, data interface{}) (err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", _UserAgent)
	resp, err := f.hc.Do(req)
	if err != nil {
		return
	}
	defer quietly.Close(resp.Body)

	jr := &_JsonResponseV1{}
	if err = json.NewDecoder(resp.Body).Decode(jr); err != nil {
		return
	}
	if jr.Ret != 0 {
		err = fmt.Errorf("error %d: %s", jr.Ret, jr.Msg)
	} else {
		err = json.Unmarshal(jr.Data, data)
	}
	return
}

var (
	errKeyMissed = errors.New("required key missed")
)

type _JsonResponseV2 struct {
	key string
	ptr interface{}
}

func (jr *_JsonResponseV2) UnmarshalJSON(data []byte) (err error) {
	dict := make(map[string]json.RawMessage)
	if err = json.Unmarshal(data, &dict); err == nil {
		// Check ret code
		var ret int
		if err = json.Unmarshal(dict["ret"], &ret); err == nil && ret != 0 {
			err = fmt.Errorf("upstream error: %d", ret)
		}
		if err != nil {
			return
		}
		// Find and parse data
		if raw, ok := dict[jr.key]; !ok {
			err = errKeyMissed
		} else {
			err = json.Unmarshal(raw, jr.ptr)
		}
	}
	return
}

func (f *fetcherImpl) getJsonV2(url string, key string, data interface{}) (err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", _UserAgent)
	resp, err := f.hc.Do(req)
	if err != nil {
		return
	}
	defer quietly.Close(resp.Body)
	// Parse response
	jr := &_JsonResponseV2{key: key, ptr: data}
	return json.NewDecoder(resp.Body).Decode(jr)
}
