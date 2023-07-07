package workload

import (
	"func/cachelet"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

const (
	ROUND = 20
)

// Workload represent the user application as function
type Workload interface {
	Run()
	ReportOutCome()
}

type BasicWorkload struct {
	store *cachelet.Cachelet
}

func NewBasicWorkload() *BasicWorkload {
	return &BasicWorkload{
		store: cachelet.New(),
	}
}

func (b *BasicWorkload) Run() {
	log.Infof("Simulation Workload Start")
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
