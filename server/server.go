package server

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/apprentice3d/forge-api-go-client/oauth"
)

//StartServer is responsible for setting up and lunching a simple web-server on available port
func StartServer() {

	service := ForgeServices{
		oauth: setupForgeOAuth(),
	}

	//serving static files
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", fs)

	// routes
	http.HandleFunc("/gettoken", service.getToken)
	http.HandleFunc("/geturn", service.getURN)
	http.HandleFunc("/upload", service.uploadFiles)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("Could not get a port: ", err)
	}


	log.Printf("Serving on port %d\n\n ", listener.Addr().(*net.TCPAddr).Port)
	if err := http.Serve(listener, nil); err != nil {
		log.Fatalln(err.Error())
	}
}

func setupForgeOAuth() oauth.TwoLeggedAuth {
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		log.Fatal("The FORGE_CLIENT_ID and FORGE_CLIENT_SECRET env vars are not set. \nExiting ...")
	}

	log.Printf("Starting app with FORGE_CLIENT_ID = %s\n", clientID)
	return oauth.NewTwoLeggedClient(clientID, clientSecret)
}
