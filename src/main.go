package main

import (
	"log"
	"net/http"
	"osoba/auth"
	"osoba/deploy"
	"osoba/webhook"
)

func main() {
	authConfig, err := auth.InitConfig(auth.Config{
		LoginFormURI:  "/login?backTo=/osoba",
		VerifyKeyFile: "/jwt-secret/secret.key",
	})
	if err != nil {
		log.Panic(err)
	}

	http.Handle("/", loggingHandler(checkMethodHandler(http.MethodGet, authHandler(authConfig, http.HandlerFunc(mainHandler)))))

	webhoocConfigs := webhook.InitConfigs([]webhook.Config{
		webhook.Config{
			&deploy.Config{
				Path:       "/aaa",
				RootPath:   "/www/html",
				ReleaseURL: "https://github.com/w-haibara/portfolio/releases/download/v1.0.8/portfolio.zip",
			},
		},
	})
	webhooksManage(webhoocConfigs)

	http.ListenAndServe(":8080", nil)
}
