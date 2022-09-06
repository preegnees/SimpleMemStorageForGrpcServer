package storage

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	m "simpleMemStor/pkg/model"
	"testing"
)

var stor m.IStreamStorage

func TestMain(m *testing.M) {

	fmt.Println("Инициализация")
	stor = NewStorage()
	fmt.Println("Сторадж инициализирван:", stor)
	exitEval := m.Run()
	fmt.Println("Остановка")
	os.Exit(exitEval)
}

func TestSavePeersAndGetPeers(t *testing.T) {
	var idCh1 m.IdChannel = "1"
	var idCh2 m.IdChannel = "2"

	data := []struct {
		Name string
		Peer m.Peer
	}{
		{"1", m.Peer{IdChannel: idCh1, Name: "N1", GrpcStream: "test1"}},
		{"2", m.Peer{IdChannel: idCh1, Name: "N2", GrpcStream: "test2"}},
		{"3", m.Peer{IdChannel: idCh1, Name: "N3", GrpcStream: "test3"}},
		{"4", m.Peer{IdChannel: idCh2, Name: "N4", GrpcStream: "test4"}},
	}

	for _, d := range data {
		stor.SavePeer(d.Peer)
	}

	peers := stor.GetPeers(idCh1)

	if len(peers) != 3 {
		t.Error("Размер мапы не соответсвует значениб", len(peers))
	}
}

func TestDeletePeerWithoutErr(t *testing.T) {
	var idCh m.IdChannel = "1"

	data := []struct {
		Name string
		Peer m.Peer
	}{
		{"1", m.Peer{IdChannel: idCh, Name: "N1", GrpcStream: "test1"}},
		{"2", m.Peer{IdChannel: idCh, Name: "N2", GrpcStream: "test2"}},
		{"3", m.Peer{IdChannel: idCh, Name: "N3", GrpcStream: "test3"}},
		{"4", m.Peer{IdChannel: idCh, Name: "N4", GrpcStream: "test4"}},
	}

	for _, d := range data {
		stor.SavePeer(d.Peer)
	}

	randPeerFromData := rand.Intn(len(data))

	err := stor.DeletePeer(data[randPeerFromData].Peer)
	if err != nil {
		t.Error("Ошибка должна равняться nil")
	}

	if len(stor.GetPeers(idCh)) != 3 {
		t.Error("колличество элементов должно равняться трем")
	}
}

func TestDeletePeerWithErr(t *testing.T) {
	var idCh m.IdChannel = "1"

	data := []struct {
		Name string
		Peer m.Peer
	}{
		{"1", m.Peer{IdChannel: idCh, Name: "N1", GrpcStream: "test1"}},
		{"2", m.Peer{IdChannel: idCh, Name: "N2", GrpcStream: "test2"}},
		{"3", m.Peer{IdChannel: idCh, Name: "N3", GrpcStream: "test3"}},
		{"4", m.Peer{IdChannel: idCh, Name: "N4", GrpcStream: "test4"}},
	}

	for _, d := range data {
		stor.SavePeer(d.Peer)
	}

	err := stor.DeletePeer(m.Peer{IdChannel: "100", Name: "N1", GrpcStream: "test1"})
	if !errors.Is(err, ErrInvalidIdChannelWhenRemove) {
		t.Error("Ошибка не совпадает с ошибкой при невалидном канале", err)
	}

	err = stor.DeletePeer(m.Peer{IdChannel: idCh, Name: "N100", GrpcStream: "test100"})
	if !errors.Is(err, ErrInvalidPeerWhenRemove) {
		t.Error("Ошибка не совпадает с ошибкой при невалидном пире", err)
	}

	if len(stor.GetPeers(idCh)) != 4 {
		t.Error("колличество элементов должно равняться четырем")
	}
}
