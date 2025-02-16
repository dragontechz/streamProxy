package main

import (
	"fmt"
	"io"

	//	"io"
	"log"
	"net"
	"strings"
)

func logRequest(method, url, addr string) {
	log.Printf("Received request: %s %s from %s", method, url, addr)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Lire la requête du client
	requestBuffer := make([]byte, 4096)
	n, err := conn.Read(requestBuffer)
	if err != nil {
		log.Printf("Error reading from connection: %v", err)
		return
	}

	// Convertir la requête en chaîne de caractères
	request := string(requestBuffer[:n])

	// Extraire la ligne de la requête (première ligne)
	var method, url string
	_, err = fmt.Sscanf(request, "%s %s", &method, &url)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		return
	}

	logRequest(method, url, conn.RemoteAddr().String())

	if method == "CONNECT" {
		var targetHost string
		_, err = fmt.Sscanf(url, "%s", &targetHost)
		if err != nil {
			log.Printf("Error parsing target host: %v", err)
			return
		}

		// Établir une connexion avec le serveur cible
		targetConn, err := net.Dial("tcp", targetHost) // Utilise l'hôte et le port fournis
		if err != nil {
			log.Printf("Could not connect to target: %v", err)
			return
		}
		defer targetConn.Close()

		// Répondre au client que la connexion est établie
		conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

		// Transférer les données entre le client et le serveur cible
		go io.Copy(targetConn, conn) // Transférer du client vers le serveur
		io.Copy(conn, targetConn)    // Transférer du serveur vers le client
	}

	//method GET etc
	// Extraire l'hôte de l'URL
	// On utilise seulement le nom d'hôte sans le chemin ni les paramètres
	host := url[len("http://"):]

	// Trouver le premier '/' pour couper le chemin
	if pathIdx := strings.Index(host, "/"); pathIdx != -1 {
		host = host[:pathIdx] // Prendre uniquement l'hôte
	}

	// Établir une connexion avec le serveur cible
	if strings.Index(host, ":") == -1 {
		host = host + ":80"
	}
	targetConn, err := net.Dial("tcp", host) // Ajouter le port par défaut si nécessaire
	//	targetConn, err := net.Dial("tcp", host+":80")
	if err != nil {
		log.Printf("Could not connect to target: %v", err)
		return
	}
	defer targetConn.Close()

	// Transférer la requête au serveur cible
	_, err = targetConn.Write(requestBuffer[:n])
	if err != nil {
		log.Printf("Error writing to target connection: %v", err)
		return
	}

	go io.Copy(conn, targetConn)
	io.Copy(targetConn, conn)
}
func main() {
	log.Println("Proxy HTTP transparent en cours d'exécution sur :8080")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		fmt.Println("new connection initialised")
		go handleConnection(conn) // Gérer chaque connexion dans une goroutine
	}
}
