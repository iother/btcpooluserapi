package main

import (
	"github.com/iother/btcpooluserapi/api"
	"github.com/iother/btcpooluserapi/config"
	"github.com/iother/btcpooluserapi/db"
	"github.com/iother/btcpooluserapi/log"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
)

var (
	Version = "20190428"
	AppName = "btcpoolapi"
)

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Version = Version

	var debugLogging bool
	var bindAddress, configPath string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "bind, b",
			Usage:       "Pool public api bind address",
			Value:       ":9527",
			Destination: &bindAddress,
		},
		cli.StringFlag{
			Name:        "config,c",
			Value:       "./config.toml",
			Usage:       "config ",
			Destination: &configPath,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.InitLogger(debugLogging)
		log.Info(app.Name, "-", app.Version)
		return nil
	}

	app.Action = func(c *cli.Context) {
		// Print a startup message.
		log.Info("Loading...")

		// Setting Config
		cfg := &config.Config{}
		if err := config.New(cfg, configPath); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		dbClient, err := db.NewMySQL(&cfg.Db)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		//defer dbsClient.CloseDB()
		log.Info("mysql db init ...")

		// Create  the server

		vApi := api.NewAPI(dbClient, cfg, Version)

		// stop the server if a kill signal is caught
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, os.Kill)
		go func() {
			<-sigChan
			dbClient.CloseDB()
			log.Info("kill signal is caught...")
			os.Exit(1)
		}()

		http.HandleFunc("/version", vApi.VersionHandler)

		http.HandleFunc("/GetUserList", vApi.GetUserList)

		listenErr := http.ListenAndServe(bindAddress, nil)

		if listenErr != nil {
			log.Fatal("Error listening on", bindAddress, listenErr)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
