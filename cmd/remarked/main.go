package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/zerok/remarked/internal/commandchain"
	"github.com/zerok/remarked/internal/config"
	"github.com/zerok/remarked/internal/token"
)

const defaultRemarkJS = "https://remarkjs.com/downloads/remark-latest.min.js"
const defaultMarkdownFile = "slides.md"
const defaultConfigFile = "remarked.yml"
const stylesheetMountPoint = "/style/_.css"

var commit, date, version string

type context struct {
	Source        string
	RemarkJS      string
	StyleSheetURL string
	Title         string
	IsGuide       bool
	IsGuided      bool
	Token         string
}

func main() {
	var configPath string
	var markdownFile string
	var templateFile string
	var addr string
	var verbose bool
	var remarkJS string
	var styleSheet string
	var title string
	var err error
	var guide bool
	var staticFolder string
	var tkn string
	var initialize bool
	var showVersion bool
	pflag.StringVar(&configPath, "config", "remarked.yml", "Path to a configuration file")
	pflag.StringVar(&title, "title", "", "Presentation title")
	pflag.StringVar(&markdownFile, "markdown-file", "", "Path to a markdown file")
	pflag.StringVar(&templateFile, "template-file", "", "Path to a template file to override the default HTML output")
	pflag.StringVar(&remarkJS, "remarkjs", "", "URL or filepath of the remark.js file")
	pflag.StringVar(&addr, "http-addr", "localhost:8000", "Start HTTP server on this address")
	pflag.StringVar(&styleSheet, "stylesheet", "", "URL or filepath of a stylesheet")
	pflag.StringVar(&staticFolder, "static-folder", "", "Path to a folder that should be served through /static")
	pflag.BoolVar(&verbose, "verbose", false, "Verbose logging")
	pflag.BoolVar(&guide, "guide", false, "Allow guided mode")
	pflag.StringVar(&tkn, "guide-token", "", "Token required for acting as guide")
	pflag.BoolVar(&initialize, "init", false, "Initialize a remarked project in the current folder")
	pflag.BoolVar(&showVersion, "version", false, "Show version information")
	pflag.Parse()

	if showVersion {
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", version, commit, date)
		os.Exit(0)
	}

	log := logrus.New()
	if verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	if initialize {
		if err := doInit(log); err != nil {
			log.WithError(err).Fatal("Failed to initialize a new remarked project in the current folder")
		}
		return
	}

	cfg, err := config.LoadFromPath(configPath)
	if err != nil {
		log.WithError(err).Fatalf("Failed to read config file: %s", configPath)
	}
	if markdownFile != "" {
		cfg.MarkdownFile = markdownFile
	}
	if templateFile != "" {
		cfg.TemplateFile = templateFile
	}
	if title != "" {
		cfg.Title = title
	}
	if remarkJS != "" {
		cfg.RemarkJS = remarkJS
	}
	if cfg.RemarkJS == "" {
		cfg.RemarkJS = defaultRemarkJS
	}
	if cfg.MarkdownFile == "" {
		log.Infof("No markdown file specified. Using %s", defaultMarkdownFile)
		cfg.MarkdownFile = defaultMarkdownFile
	}
	if styleSheet != "" {
		cfg.Stylesheet = styleSheet
	}
	if guide {
		if tkn == "" {
			tkn = token.Generate()
		}
		cfg.Token = tkn
	}

	if staticFolder != "" {
		cfg.StaticFolder = staticFolder
	}

	if guide {
		log.Infof("Starting guide mode with this token:\n\n  %s\n\n", cfg.Token)
	}

	if cfg.Title == "" {
		log.Info("No title specified. Using the name of the containing folder instead.")
		cfg.Title, err = getFolderName(markdownFile)
		if err != nil {
			log.WithError(err).Fatalf("Failed to determine name of %s's parent folder", markdownFile)
		}
	}

	srv := http.Server{
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
	}
	srv.Addr = addr
	mux := http.NewServeMux()
	srv.Handler = mux

	funcs := templateFuncs{}

	localStylesheet, ok := isLocalStylesheet(cfg.Stylesheet)
	if ok {
		mux.HandleFunc(stylesheetMountPoint, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, localStylesheet)
		})
		cfg.FinalStylesheet = stylesheetMountPoint
	} else if cfg.Stylesheet != "" {
		cfg.FinalStylesheet = cfg.Stylesheet
	}

	if guide {
		hub := commandchain.Hub{Log: log}
		mux.HandleFunc("/guide/login", guideLoginHandler(cfg, log))
		mux.HandleFunc("/guide", token.Require(cfg.Token, "/guide/login", guideHandler(cfg, log)))
		mux.HandleFunc("/ws/guide", guideWebsocketHandler(cfg, &hub, log))
		mux.HandleFunc("/ws/guided", guidedWebsocketHandler(cfg, &hub, log))
	}

	if cfg.StaticFolder != "" {
		fullStaticFolder, err := filepath.Abs(cfg.StaticFolder)
		if err != nil {
			log.WithError(err).Fatalf("Failed to resolve absolute path to static folder %s", cfg.StaticFolder)
		}
		log.Debugf("Serving static files from %s", fullStaticFolder)
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(fullStaticFolder))))
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := loadOutputTemplate(cfg.TemplateFile)
		if err != nil {
			log.WithError(err).Errorf("Failed to parse template ")
			http.Error(w, "Failed to parse template file", http.StatusInternalServerError)
			return
		}
		data, err := ioutil.ReadFile(cfg.MarkdownFile)
		if err != nil {
			log.WithError(err).Errorf("Failed to read %s", cfg.MarkdownFile)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		ctx := &context{
			RemarkJS:      cfg.RemarkJS,
			StyleSheetURL: cfg.FinalStylesheet,
			Title:         cfg.Title,
			IsGuided:      guide,
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
	})

	log.Infof("Starting server on %s", addr)
	log.Debugf("Final configuration: %s", cfg)
	if err := srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatalf("Failed to start server on %s", addr)
	}
}

func buildContent(rawContent string, cfg *config.Config, fmap template.FuncMap) (string, error) {
	if cfg.MarkdownAsTemplate {
		var content bytes.Buffer
		contentTmpl, err := template.New("content").
			Funcs(fmap).
			Delims(cfg.LeftActionDelimiter, cfg.RightActionDelimiter).
			Parse(rawContent)
		if err != nil {
			return "", fmt.Errorf("Failed to parse file: %s", err.Error())
		}
		if err := contentTmpl.Execute(&content, context{}); err != nil {
			return "", fmt.Errorf("Failed to render template: %s", err.Error())
		}
		return content.String(), nil
	}
	return rawContent, nil
}

func isLocalStylesheet(u string) (string, bool) {
	if u == "" {
		return "", false
	}
	if _, err := os.Stat(u); err != nil {
		return "", false
	}
	return u, true
}

func getFolderName(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Base(filepath.Dir(abs)), nil
}

func doInit(log *logrus.Logger) error {
	cfgFileStat, err := os.Stat(defaultConfigFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to check %s: %s", defaultConfigFile, err.Error())
		}
	} else {
		if cfgFileStat.IsDir() {
			return fmt.Errorf("%s exists but is a directory", defaultConfigFile)
		}
	}
	cfgDir := filepath.Dir(defaultConfigFile)
	if err := os.MkdirAll(cfgDir, 0700); err != nil {
		return fmt.Errorf("%s could not be created: %s", cfgDir, err.Error())
	}
	log.Infof("Generating empty config file: %s", defaultConfigFile)
	if err := ioutil.WriteFile(defaultConfigFile, []byte(config.Sample), 0644); err != nil {
		return fmt.Errorf("failed to create %s: %s", defaultConfigFile, err.Error())
	}
	return nil
}
