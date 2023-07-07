package cachelet

import (
	"func/memberdaemon"
	"func/util"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	PeerStatusConnected = "PeerStatusConnected"
	PeerStatusUnknown   = "PeerStatusUnknown"
)

type Cachelet struct {
	data                       map[string]interface{}
	remoteStorageFetchDuration int
	metadata                   *Metadata
	peers                      map[string]*Peer
	md                         memberdaemon.MemberDemon
	mu                         sync.Mutex
}

type Metadata struct {
	localHitNum   int
	localMissNum  int
	remoteHitNum  int
	remoteMissNum int
}

type Peer struct {
	endpoint string
	status   string
}

type GetResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var onceCachelet sync.Once
var cachelet *Cachelet

func New() *Cachelet {
	onceCachelet.Do(func() {
		log.Infof("Once new cachelet")
		cachelet = &Cachelet{
			data:                       make(map[string]interface{}),
			remoteStorageFetchDuration: 600,
			peers:                      make(map[string]*Peer),
			metadata:                   &Metadata{},
		}

		selfEndpoint, peerEndpointsStr := "", ""

		// used in non-container dev
		if len(os.Args) == 3 {
			os.Setenv("SELF_ENDPOINT", os.Args[1])
			os.Setenv("PEER_ENDPOINTS", os.Args[2])

		}

		util.MustMapEnv("PEER_ENDPOINTS", &peerEndpointsStr)
		util.MustMapEnv("SELF_ENDPOINT", &selfEndpoint)

		peerEndpoints := strings.Split(peerEndpointsStr, ",")

		for _, endpoint := range peerEndpoints {
			cachelet.peers[endpoint] = &Peer{
				endpoint: endpoint,
				status:   PeerStatusUnknown,
			}
		}
	})
	return cachelet
}

func (c *Cachelet) Init() error {
	go func() {
		allConnected := false
		for !allConnected {
			allConnected = true
			for _, peer := range c.peers {
				allConnected = allConnected && c.connect(peer)
			}
		}
		time.Sleep(2 * time.Second)
	}()
	c.initServer()
	return nil
}

func (c *Cachelet) connect(peer *Peer) bool {
	if ok := sendConnectionRequest(peer.endpoint); ok {
		c.handleConnection(peer.endpoint)
		return true
	} else {
		return false
	}
}

func (c *Cachelet) handleConnection(endpoint string) error {
	c.mu.Lock()
	if _, ok := c.peers[endpoint]; ok {
		c.peers[endpoint].status = PeerStatusConnected
	} else {
		c.peers[endpoint] = &Peer{endpoint: endpoint, status: PeerStatusConnected}
	}
	c.md.AddPeer(endpoint)
	c.mu.Unlock()
	return nil
}

func (c *Cachelet) Get(key string) interface{} {
	value, exist := c.data[key]
	if exist {
		c.metadata.localHitNum++
		return value
	}

	log.Infof("[Local Miss] Data for key %s not in local,  get from remote storage", key)
	//value = c.getFromPeers(key)

	//if value != nil {
	//	c.metadata.remoteHitNum++
	//	return value
	//}
	//c.metadata.remoteMissNum++
	value = c.getFromRemoteStorage(key)
	c.data[key] = value

	return value
}

// This function is for peer to get data from it'c local cachelet
func (c *Cachelet) PeersGet(key string) interface{} {
	value, exist := c.data[key]
	if exist {
		return value
	}
	return nil
}

// TODO: maybe really fetch from remote storage?
func (c *Cachelet) getFromRemoteStorage(key string) interface{} {
	time.Sleep(time.Millisecond * time.Duration(c.remoteStorageFetchDuration))
	k, err := strconv.Atoi(key)
	if err != nil {
		log.Error(err.Error())
	}
	v := struct{ name string }{name: "remote storage"}
	log.Infof("[Fetching][remote Storage] Fetching from remote storage key: %d, return value: %d", k, v)

	return v
}

func (c *Cachelet) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *Cachelet) Report() {
	log.Printf("Local hit: %d\n", c.metadata.localHitNum)
	log.Printf("Local miss: %d", c.metadata.localHitNum)
	log.Printf("Remote hit: %d\n", c.metadata.remoteHitNum)
	log.Printf("Remote miss: %d", c.metadata.remoteMissNum)
}
