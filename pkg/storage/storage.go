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
	streams              map[m.IdChannel](map[m.Peer]struct{})
	chs                  map[m.IdChannel](chan map[m.Peer]struct{})
	NotifyOfNewConn      chan []m.Peer
	CloseNotifyOfNewConn func()
	mx                   sync.Mutex
}

// Функция получения хранилища
func NewStorage() m.IStreamStorage {
	strg := make(map[m.IdChannel](map[m.Peer]struct{}))
	cs := make(map[m.IdChannel](chan map[m.Peer]struct{}))
	notify := make(chan []m.Peer)
	funcCloseNotify := func() {
		close(notify)
	}
	return &storage{
		streams:              strg,
		chs:                  cs,
		NotifyOfNewConn:      notify,
		CloseNotifyOfNewConn: funcCloseNotify,
	}
}

// Сохранение пира при подключении
func (s *storage) SavePeer(peer m.Peer) { // нужно возвращать канал, если его нет то создать

	s.mx.Lock()
	defer s.mx.Unlock()

	peers, _ := s.streams[peer.IdChannel]

	if peers == nil {
		peers = make(map[m.Peer]struct{})
	}

	peers[peer] = struct{}{}
	s.streams[peer.IdChannel] = peers
	return
}

// Удаление пира при отключении
func (s *storage) DeletePeer(peer m.Peer) error {

	s.mx.Lock()
	defer s.mx.Unlock()

	peers, ok := s.streams[peer.IdChannel]
	if !ok {
		return ErrInvalidIdChannelWhenRemove
	}

	_, ok = peers[peer]
	if !ok {
		return ErrInvalidPeerWhenRemove
	}

	delete(peers, peer)
	s.streams[peer.IdChannel] = peers
	return nil
}

// Получение всех пиров
func (s *storage) GetPeers(idCh m.IdChannel) map[m.Peer]struct{} {

	s.mx.Lock()
	defer s.mx.Unlock()

	peers, _ := s.streams[idCh]
	return peers
}

// Для каждого idchannel будет создаваться свой канал, куда будет писать этот писарь
func (s *storage) sendPeers(idCh m.IdChannel) {
	go func() {
		ch, ok := s.chs[idCh]
		if ok {
			peers := s.GetPeers(idCh)
			if len(peers) >= 2 {
				ch <- peers
			}
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
