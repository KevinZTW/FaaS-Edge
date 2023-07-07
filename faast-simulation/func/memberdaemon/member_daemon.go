package memberdaemon

import "github.com/lafikl/consistent"

type MemberDemon struct {
	peers *consistent.Consistent
}

func New(endpoint string) *MemberDemon {
	peers := consistent.New()
	peers.Add(endpoint)

	return &MemberDemon{
		peers: peers,
	}
}

func (m *MemberDemon) GetOwner(key string) string {
	if owner, err := m.peers.Get(key); err != nil {
		panic(err.Error())
	} else {
		return owner
	}
}

func (m *MemberDemon) AddPeer(endpoint string) {
	m.peers.Add(endpoint)
}

func (m *MemberDemon) GetPeerEndpoints() []string {
	return m.peers.Hosts()
}
