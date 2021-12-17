package handlers

import (
	"fmt"
	"log"
	"net/http"
	"ws/internal/dto"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan dto.WsPayload)
var clients = make(map[dto.WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.htm", nil)
	if err != nil {
		log.Println(err)

	}
}

// WsEndPoint upgrades connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to endpoint")
	response := dto.WsJsonResponse{
		Message: `<em><small>Connected to server</small></em>`,
	}
	conn := dto.WebSocketConnection{Conn: ws}
	clients[conn] = ""
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
	go ListenForWs(&conn)
}

func ListenForWs(conn *dto.WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()
	payload := dto.WsPayload{}
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func getUserList() []string {
	userList := []string{}
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	return userList
}

func ListenToWsChannel() {
	response := dto.WsJsonResponse{}
	for {
		e := <-wsChan
		switch e.Action {
		case "username":
			clients[e.Conn] = e.Username
			userList := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = userList
			broadcastToAll(response)
			//get a list of all users and send it back via Broadcast
			break
		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)
			response.ConnectedUsers = getUserList()
			broadcastToAll(response)
			break
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf(`<strong>%s</strong>: %s`, e.Username, e.Message)
			broadcastToAll(response)
		}

	}
}

func broadcastToAll(response dto.WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
