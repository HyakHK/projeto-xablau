package main

import(
	"log"
	"net/http"
	"sync"
	
	"github.com/gorilla/websocket"
)

var (
	//upgrade to persistent websocket
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
)


//continuar a partir daqui, pausa
var (
	ErrEventNotSupported = errors.New("Não há suporte a este tipo de evento")
)

//all registers references

type Manager struct {
	clients ClientList

	//lcok state before diting
	sync.RWMutex

	//handle event functions
	handlers map[string]EventHandler

}

//initialize the values inside manager
func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

//serveWS Http handler that has manager that allows connection
func (m * Manager) serveWS(w http.ResponseWriter, r *http.Request){
	log.Println("New connection")
	//upgrade HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//create new client
	client := NewClient(conn, m)
	//add to manager
	m.addClient(client)
	//start write read message process
	go client.readMessages()
}

//add client to list
func(m *Manager) addClient(client *Client){
	//manipulate
	m.Lock()
	defer m.Unlock()

	//add client
	m.clients[client]= true
}

//remove client
func (m *Manager) removeClient(client *Client){
	m.Lock()
	defer m.Unlock()

	//if exists delete
	if _, ok :=m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}