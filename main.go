package main

import (
	"drtech.co/gl2gl/core"
	"drtech.co/gl2gl/core/configs"
	"drtech.co/gl2gl/orm"
	"drtech.co/gl2gl/services"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()
	app := &cli.App{
		Name:  "gl2gl",
		Usage: "sync gitlab to gitlab",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "db",
				Value:   "./assets/db",
				Usage:   "-db ./assets/db",
				EnvVars: []string{"DB"},
			},
			&cli.StringFlag{
				Name:    "ll",
				Value:   "trace",
				Usage:   "-ll trace ",
				EnvVars: []string{"LOGLEVEL"},
			},
			&cli.StringFlag{
				Name:    "td",
				Value:   "./tmp/gl2gl",
				Usage:   "-td /tmp/gl2gl ",
				EnvVars: []string{"TD"},
			},
		},
		Action: func(c *cli.Context) error {
			db := c.String("db")
			td := c.String("td")
			logLevelStr := c.String("ll")
			logLevel, err := logrus.ParseLevel(logLevelStr)
			if err != nil {
				return err
			}
			logrus.SetLevel(logLevel)
			logger := logrus.WithField("Name", "main")
			logger.Info("PublishTime:", core.PublishTime)
			logger.Info("VERSION:", core.VERSION)
			logger.Info("DB:", db)
			logger.Info("TD:", td)
			configs.SqliteDsn = fmt.Sprintf("file:%s?cache=shared", db)
			configs.TempDir = td
			err = orm.Setup()
			if err != nil {
				return err
			}
			err = services.Setup()
			if err != nil {
				return err
			}
			for {
				c := make(chan os.Signal)
				signal.Notify(c, os.Interrupt, syscall.SIGTERM)
				<-c
				os.Exit(1)
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
