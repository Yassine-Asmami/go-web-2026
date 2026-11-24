package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var JeuxNoms = map[int]string{
	1: "Celeste",
	2: "Hades",
	3: "Rocket League",
	4: "Street Fighter 6",
	5: "Stardew Valley",
}

func main() {
	port := flag.String("port", "8080", "port d'écoute")
	flag.Parse()

	games, err := FetchGames()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur :", err)
		os.Exit(1)
	}
	tournaments, err := FetchTournaments()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur :", err)
		os.Exit(1)
	}
	printJeux(games, tournaments)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", Accueil)
	mux.HandleFunc("GET /jeux", listJeux)
	mux.HandleFunc("GET /jeux/{id}", Jeux)

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	addr := ":" + *port
	log.Printf("GameDex sur http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func Accueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
	<html lang="fr">
	<head>
  		<meta charset="utf-8">
  		<title>GameDex</title>
  		<link rel="stylesheet" href="/static/css/style.css">
	</head>
	<body>
  		<h1>GameDex</h1>
  		<p>Le catalogue de jeux vidéo et de leurs tournois.</p>
  		<p><a href="/jeux">Voir les jeux</a></p>
	</body>
	</html>`)
}

func listJeux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Liste des jeux")
	for id := 1; id <= len(JeuxNoms); id++ {
		fmt.Fprintf(w, "%d. %s\n", id, JeuxNoms[id])
	}
}

func Jeux(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "l'identifiant doit être un nombre", http.StatusBadRequest)
		return
	}
	name, ok := JeuxNoms[id]
	if !ok {
		http.Error(w, "jeu introuvable", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Jeu %d : %s\n", id, name)
}

func printJeux(games []Game, tournaments []Tournament) {
	count := map[int]int{}
	for _, t := range tournaments {
		count[t.GameID]++
	}
	for _, g := range games {
		fmt.Printf("%s (%d) : %d tournoi(s)\n", g.Name, g.Year, count[g.ID])
	}
}
