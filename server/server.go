package server

import (
	"log"
	"net/http"
	"os"

	"github.com/apprentice3d/forge-api-go-client/oauth"
)

//StartServer is responsible for setting up and lunching a simple web-server on the specified port
func StartServer(port string) {

	service := ForgeServices{
		oauth: setupForgeOAuth(),
	}

	//serving static files
	fs := http.FileServer(http.Dir("client/build"))
	http.Handle("/", fs)

	// routes
	http.HandleFunc("/gettoken", service.getToken)
	http.HandleFunc("/geturn", service.getURN)
	http.HandleFunc("/upload", service.uploadFiles)

	log.Printf("Serving on port %s\n\n ", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalln(err.Error())
	}
}

func setupForgeOAuth() oauth.AuthApi {
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		log.Fatal("The FORGE_CLIENT_ID and FORGE_CLIENT_SECRET env vars are not set. \nExiting ...")
	}

	log.Printf("Starting app with FORGE_CLIENT_ID = %s\n", clientID)
	return oauth.NewTwoLeggedClient(clientID, clientSecret)
}
