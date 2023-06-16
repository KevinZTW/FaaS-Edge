package app

import (
	"func/store"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

const ROUND = 20

type Workload interface {
	Run()
	ReportOutCome()
}

type BasicWorkload struct {
	store *store.Store
}

func NewBasicWorkload(peerAddress string) *BasicWorkload {
	return &BasicWorkload{
		store: store.New(peerAddress),
	}
}

func (b *BasicWorkload) Run() {
	randSource := rand.NewSource(time.Now().UnixNano())
	ran := rand.New(randSource)

	//randomly check user's ID (int) to get their account balance
	for i := 0; i < ROUND; i++ {
		key := strconv.Itoa(ran.Intn(10))
		value := b.store.Get(key)
		log.Infof("user id:%s got balance %d", key, value)
	}
}

func (b *BasicWorkload) ReportOutCome() {
	b.store.Report()
}
