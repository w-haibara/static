package osoba

import (
	"github.com/k0kubun/pp"
)

type App struct {
	DocumentRoot         string
	TmpDirContentsPrefix string
	Contents             Contents
}

func Run() {
	a := App{}
	a.LoadConfig()
	pp.Println(a)

	a.FetchAll()
}
