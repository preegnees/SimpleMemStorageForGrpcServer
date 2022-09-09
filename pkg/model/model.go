package model

// IdChannel - айди канала, с которы свзяаны пиры
type IdChannel string

/*
Peer. Пир - это подключение, которое имеет:
Name=Имя(клиента);
Connected=статус подключения, если false, то стрим скорее всего nil;
GrpcStream=стрим gprc;
*/
type Peer struct {
	IdChannel  IdChannel
	Name       string
	GrpcStream interface{}
}

/*
IMemStorage. Интерфейс для взаимодействия с базой данных.
SavePeer=сохранение пира в соответвии с его id канала;
DeletePeer=удалять пир из пула, если он отключился;
*/
type IStreamStorage interface {
	SavePeer(peer Peer) <-chan map[Peer]struct{}
	DeletePeer(peer Peer) error
	GetPeers(idCh IdChannel) map[Peer](chan map[Peer]struct{})
}
