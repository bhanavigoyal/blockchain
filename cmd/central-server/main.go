package centralserver

import (
	"log"
	"net/http"
)


func main(){
	setupAPI()
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080",nil))
}

func setupAPI(){
	manager := NewManager()
	http.HandleFunc("/ws", manager.serveWs)
}