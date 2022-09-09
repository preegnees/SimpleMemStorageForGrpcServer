package storage

import (
	"fmt"
	"sync"

	m "simpleMemStor/pkg/model"
)

var ErrInvalidIdChannelWhenRemove error = fmt.Errorf("Ошибка при удалении пира, такого канала не существует")
var ErrInvalidPeerWhenRemove error = fmt.Errorf("Ошибка при удалении пира, такого пира не существует")

// Проверка на соответсвии интерфейсу
var _ m.IStreamStorage = (*storage)(nil)

// Структура хранилища, состоит из мапы из айди канала и пиров этого канала
type storage struct {
	streams map[m.IdChannel](map[m.Peer](chan map[m.Peer]struct{}))
	mx      sync.Mutex
}

// Функция получения хранилища
func NewStorage() m.IStreamStorage {
	strg := make(map[m.IdChannel](map[m.Peer](chan map[m.Peer]struct{})))
	return &storage{
		streams: strg,
	}
}

// Сохранение пира при подключении
func (s *storage) SavePeer(peer m.Peer) <-chan map[m.Peer]struct{} { // нужно возвращать канал, если его нет то создать

	s.mx.Lock()
	defer s.mx.Unlock()

	peers, ok := s.streams[peer.IdChannel]

	if peers == nil || !ok {
		peers = make(map[m.Peer](chan map[m.Peer]struct{}))
	}

	ch := make((chan map[m.Peer]struct{}))
	peers[peer] = ch
	s.streams[peer.IdChannel] = peers

	go s.SendPeers(peer.IdChannel)

	return ch
}

// Удаление пира при отключении
func (s *storage) DeletePeer(peer m.Peer) error {

	s.mx.Lock()
	defer s.mx.Unlock()

	peers, ok := s.streams[peer.IdChannel]
	if !ok {
		return ErrInvalidIdChannelWhenRemove
	}

	ch, ok := peers[peer]
	if !ok {
		return ErrInvalidPeerWhenRemove
	}

	delete(peers, peer)
	close(ch)
	s.streams[peer.IdChannel] = peers

	go s.SendPeers(peer.IdChannel)

	return nil
}

// Для каждого idchannel будет создаваться свой канал, куда будет писать этот писарь
// по идее у нас для каждого пира есть свой канала, через который он будет что то узнавать
// эта функция должна вызываться каждый раз когда есть изменение в каком то id channel
func (s *storage) SendPeers(idCh m.IdChannel) {
	peers, ok := s.streams[idCh]
	if !ok {
		return
	} 
	ps := make(map[m.Peer]struct{}) // сохранение отдельно пиров
	chs := make(map[chan map[m.Peer]struct{}]struct{}) // отдельно каналов
	for k, v := range peers {
		ps[k] = struct{}{}
		chs[v] = struct{}{}
	}
	send(ps, chs)
}

func (s *storage) CloseChan(peer m.Peer) {
	strm, ok := s.streams[peer.IdChannel]
	if ok {
		ch := strm[peer]
		if ch != nil {
			close(ch)
		}
	}
}

func send(ps map[m.Peer]struct{}, chs map[chan map[m.Peer]struct{}]struct{}) {
	go func ()  {
		for	k := range chs {
			go func(k chan map[m.Peer]struct{}) {
				k <- ps
			}(k)
		}
	}()
}

// func (s *myStorage) SaveStream(streamNew m.Stream) (string, error) {
// 	streamOld, ok := s.streams[streamNew.IdChannel]
// 	if ok {
// 		log.Println(streamNew)
// 		log.Println(streamOld)
// 		if len(streamNew.Peers) != len(streamOld.Peers) {
// 			return "", fmt.Errorf("Одинаковый id канал, но разные значения, зачит вы где то ошиблись (разная длинна)")
// 		}
// 		count := 0
// 		l := len(streamNew.Peers)
// 		for _, pOld := range streamOld.Peers {
// 			for _, pNew := range streamNew.Peers {
// 				if pNew.Name == pOld.Name {
// 					count++
// 				}
// 			}
// 		}
// 		if l == count {
// 			return "Такой стрим уже есть", nil
// 		}
// 		return "", nil
// 	} else {
// 		s.streams[streamNew.IdChannel] = streamNew
// 		log.Println("Было сохранено:", streamNew)
// 		return "", nil
// 	}
// }

// func (s *myStorage) SavePeer(idChannel string, peer m.Peer) error {
// 	return nil
// }

// func (s *myStorage) SetFiledConnected(idChannel string, peer m.Peer, is bool) error {
// 	return nil
// }

// func (s *myStorage) GetStreams() (strm chan<- m.Stream) {
// 	return nil
// }

// func New() m.IStreamStorage {
// 	return &myStorage{
// 		streams: make(map[string]m.Stream, 0),
// 		// peers:   make(map[m.Peer]struct{}, 0),
// 	}
// }
