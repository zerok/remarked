package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/zerok/remarked/internal/commandchain"
	"github.com/zerok/remarked/internal/config"
)

func guidedWebsocketHandler(cfg *config.Config, hub *commandchain.Hub, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			HandshakeTimeout: time.Second * 2,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.WithError(err).Error("Failed to upgrade connection")
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		recv := &commandchain.Receiver{Conn: conn, Log: log}
		hub.RegisterReceiver(recv)
		defer hub.UnregisterReceiver(recv)
		recv.Handle(r.Context())
		if err := recv.Handle(r.Context()); err != nil {
			log.WithError(err).Error("Receiver exited")
		}
	}
}

func guideWebsocketHandler(cfg *config.Config, hub *commandchain.Hub, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			HandshakeTimeout: time.Second * 2,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.WithError(err).Error("Failed to upgrade connection")
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		cmdr := &commandchain.Commander{Conn: conn, Log: log, Token: cfg.Token}
		hub.RegisterCommander(cmdr)
		defer hub.UnregisterCommander(cmdr)
		if err := cmdr.Handle(r.Context()); err != nil {
			log.WithError(err).Error("Commander exited")
		}
	}
}

func guideHandler(cfg *config.Config, tmpl *template.Template, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(cfg.MarkdownFile)
		if err != nil {
			log.WithError(err).Errorf("Failed to read %s", cfg.MarkdownFile)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		funcs := templateFuncs{}
		ctx := context{
			RemarkJS:      cfg.RemarkJS,
			StyleSheetURL: cfg.FinalStylesheet,
			Title:         cfg.Title,
			IsGuided:      true,
			IsGuide:       true,
			Token:         cfg.Token,
		}

		rawData := string(data)
		content, err := buildContent(rawData, cfg, funcs.FuncMap())
		if err != nil {
			log.WithError(err).Error("Failed to compile content")
			http.Error(w, "Failed to compile output", http.StatusInternalServerError)
			return
		}
		ctx.Source = content
		tmpl.Execute(w, ctx)
	}
}

func guideLoginHandler(cfg *config.Config, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.FormValue("token") == cfg.Token {
				w.Header().Set("Set-Cookie", fmt.Sprintf("guideToken=%s", cfg.Token))
				w.Header().Set("Location", "/guide")
				w.WriteHeader(307)
				return
			}
		}
		fmt.Fprintf(w, `<!doctype html>
			<html>
				<head>
					<title>Login</title>
					<meta charset="utf-8">
					<meta name="viewport" content="width=device-width, initial-scale=1">
					<style>
					input {
						display: block;
					}
					.form__actions {
						margin-top: 5px;
					}
					</style>
				</head>
				<body>
					<form method="post">
						<label>Token: 
							<input type="password" name="token" />
						</label>
						<div class="form__actions">
							<button type="submit">Authenticate</button>
						</div>
					</form>
				</body>
			</html>`)
	}
}
