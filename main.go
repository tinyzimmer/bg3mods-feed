package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"

	"github.com/tinyzimmer/bg3mods-feed/internal/config"
	"github.com/tinyzimmer/bg3mods-feed/internal/feed"
	"github.com/tinyzimmer/bg3mods-feed/internal/mods"
	"github.com/tinyzimmer/bg3mods-feed/internal/server"
)

var (
	Version = "dev"
	Commit  = "none"
)

func main() {
	flags := pflag.NewFlagSet("bg3mods-feed", pflag.ExitOnError)
	configFile := flags.String("config", "", "Path to the configuration file (YAML, JSON, TOML, or HCL)")
	version := flags.Bool("version", false, "Print the version and exit")
	config.BindPFlags(flags)
	flags.Parse(os.Args[1:])

	if *version {
		log.Printf("bg3mods-feed %s (%s)", Version, Commit)
		os.Exit(0)
	}

	conf, err := config.Load(*configFile)
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	log.Printf("Loaded configuration: %+v", conf)

	fetcher := mods.NewFetcher(conf.APIURL)
	generator := feed.NewGenerator(fetcher, feed.GeneratorOptions{
		MaxItems:      conf.MaxFeedItems,
		Sort:          conf.Sort,
		Tags:          conf.Tags,
		Platform:      conf.Platform,
		FetchInterval: conf.FetchInterval,
		Format:        conf.Format,
	})

	server := server.NewServer(server.ServerOptions{
		Generator: generator,
		Addr:      conf.Listen,
	})

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sigc

	log.Println("Shutting down server...")
	ctx := context.Background()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutdown server:", err)
	}
}
