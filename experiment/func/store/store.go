package store

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Store struct {
	data                       map[string]interface{}
	peerAddr                   string
	remoteStorageFetchDuration int
	localHitNum                int
	localMissNum               int
	remoteHitNum               int
	remoteMissNum              int
}

type GetResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var onceStore sync.Once
var store *Store

func New(peerAddress string) *Store {
	onceStore.Do(func() {
		var addr string
		if peerAddress != "" {
			addr = peerAddress
		} else {
			addr = os.Getenv("PEER_ADDR")
		}

		if addr == "" {
			panic("Cant find peer addr from env: PEER_ADDR or argument list")
		}
		log.Infof("Init Store")
		store = &Store{
			data:                       make(map[string]interface{}),
			peerAddr:                   addr,
			remoteStorageFetchDuration: 600,
		}
	})
	return store
}

func (s *Store) Get(key string) interface{} {
	value, exist := s.data[key]
	if exist {
		s.localHitNum++
		return value
	}

	s.localMissNum++
	log.Infof("[Local Miss] Data for key %s not in local, try to get from peers", key)
	value = s.getFromPeers(key)

	if value != nil {
		s.remoteHitNum++
		return value
	}
	s.remoteMissNum++
	value = s.getFromRemoteStorage(key)
	s.data[key] = value

	return value
}

// This function is for peer to get data from it's local store
func (s *Store) PeersGet(key string) interface{} {
	value, exist := s.data[key]
	if exist {
		return value
	}
	return nil
}

func (s *Store) getFromPeers(key string) interface{} {
	if s.peerAddr == "" {
		return nil
	}
	log.Infof("[Fetching][peer] addr: %s", s.peerAddr)

	url := fmt.Sprintf("%s/store?key=%s", s.peerAddr, key)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Error(err.Error())
		return nil
	}
	gr := GetResponse{}

	log.Infof("==========%s", body)
	err = json.Unmarshal(body, &gr)

	if err != nil {
		log.Error(err.Error())
		return nil
	}
	log.Infof("==========%v", gr)

	return gr.Value
}

// TODO: maybe really fetch from remote storage?
func (s *Store) getFromRemoteStorage(key string) interface{} {
	time.Sleep(time.Millisecond * time.Duration(s.remoteStorageFetchDuration))
	k, err := strconv.Atoi(key)
	if err != nil {
		log.Error(err.Error())
	}
	v := struct{ name string }{name: "remote storage"}
	log.Infof("[Fetching][remote Storage] Fetching from remote storage key: %d, return value: %d", k, v)

	return v
}

func (s *Store) Set(key string, value interface{}) {
	s.data[key] = value
}

func (s *Store) Report() {
	log.Printf("Local hit: %d\n", s.localHitNum)
	log.Printf("Local miss: %d", s.localMissNum)
	log.Printf("Remote hit: %d\n", s.remoteHitNum)
	log.Printf("Remote miss: %d", s.remoteMissNum)
}
