package model

type Peer struct {
	Name      string
	Connected bool
	Stream    interface{}
}

type Stream struct {
	IdChannel string
	Peers     []Peer
}

// IMemStorage ...
type IStreamStorage interface {
	SaveStream(strm Stream) (string, error)
	SavePeer(idChannel string, peer Peer) error
	SetFiledConnected(idChannel string, peer Peer, is bool) error
	GetStreams() (strm chan<- Stream)
}

// IAccessStorage ...
type IAccessStorage interface {
	Check(string) (bool, error)
	SaveNewName(string) (bool, error)
}

// IStorage ...
type IStorage interface {
	IAccessStorage
	IStreamStorage
}