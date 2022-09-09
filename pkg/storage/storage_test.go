package storage

import (
	"fmt"
	"os"
	"sync"
	"time"

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
	// var idCh2 m.IdChannel = "2"

	data := []struct {
		Name string
		Peer m.Peer
	}{
		{"1", m.Peer{IdChannel: idCh1, Name: "N1", GrpcStream: "test1"}},
		{"2", m.Peer{IdChannel: idCh1, Name: "N2", GrpcStream: "test2"}},
		{"3", m.Peer{IdChannel: idCh1, Name: "N3", GrpcStream: "test3"}},
		// {"4", m.Peer{IdChannel: idCh2, Name: "N4", GrpcStream: "test4"}},
	}

	chs := make([]<-chan map[m.Peer]struct{}, 0)

	for _, d := range data {
		chs = append(chs, stor.SavePeer(d.Peer))
	}
	

	if len(chs) == 0 {
		fmt.Println("Длинна массива chs = 0")
	}

	count := 0
	var mx sync.Mutex
	incr := func () {
		mx.Lock()
		defer mx.Unlock()
		count++
	}

	var wg sync.WaitGroup
	wg.Add(len(chs) + 1)

	for i, ch := range chs {
		go func (i int, ch <-chan map[m.Peer]struct{})  {
			for {
				select {
				case val, ok := <-ch:
					if ok {
						incr()
						fmt.Println("index=", i, ", val=", val)	
						continue
					}
					wg.Done()
					fmt.Println("index=", i, ", closed")
					return
				}
			}
		}(i, ch)
	}

	time.Sleep(1 * time.Second)
	err := stor.DeletePeer(m.Peer{IdChannel: idCh1, Name: "N1", GrpcStream: "test1"})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(1 * time.Second)
	i := 5
	myPeer := m.Peer{IdChannel: idCh1, Name: "N3000", GrpcStream: "test3000"}
	ch := stor.SavePeer(myPeer)

	go func (i int, ch <-chan map[m.Peer]struct{})  {
		for {
			select {
			case val, ok := <-ch:
				if ok {
					incr()
					fmt.Println("index=", i, ", val=", val)	
					continue
				}
				wg.Done()
				fmt.Println("index=", i, ", closed")
				return
			}
		}
	}(i, ch)

	time.Sleep(4 * time.Second)
	for _, d := range data {
		stor.CloseChan(d.Peer)
	}
	stor.CloseChan(myPeer)
	wg.Wait()

	fmt.Println(count)

	if count != len(chs) {
		t.Error("не равно")
	}
}