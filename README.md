# Golang WebSocket Chat
This is a Golang Web socket chat application that can run locally and is implemented with help of Gorilla websocket and it renders frontend with help of Golang Templating Library called as Jet and it uses Reconnecting websocket javascript lib for reconnecting broken connections

To run it , cd cmd/web && go run ./


Now open the app on localhost:8080 and open the app in different tabs , and type in the username and message , each tab opens a new websocket and you have an example of a live chat room with help of websockets
