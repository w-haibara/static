package osoba

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/k0kubun/pp"
)

type App struct {
	DocumentRoot         string
	TmpDirContentsPrefix string
	Contents             Contents
}
type Path string
type URL string
type Contents struct {
	Mu sync.Mutex
	V  map[Path]URL
}

func Run() {
	a := App{}
	a.LoadConfig()
	pp.Println(a)

	chanDeployPath := make(chan Path)

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/api/content", helloHandler)

	go func() {
		for {
			a.Deploy(<-chanDeployPath)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from a HandleFunc #1!\n")
}
