package model

// IdChannel - канала, который связывет пиры
type IdChannel string

// Name - имя клиента
type Name string

// Token - токен для связки в хранилище каналы и пиры
type Token string

/*
Peer. Пир - это подключение, которое имеет:
IdChannel: id канала, с которым он связан;
Name: имя, которое имеет пир;
GrpcStream=стрим gprc;
*/
type Peer struct {
	IdChannel  IdChannel
	Name       Name
	GrpcStream interface{}
}

/*
IMemStorage. Интерфейс для взаимодействия с базой данных.
SavePeer: сохранение пира;
DeletePeer: удаление пира при отключении;
*/
type IStreamStorage interface {
	SavePeer(peer Peer) <-chan map[Peer]struct{}
	DeletePeer(peer Peer) error
}
