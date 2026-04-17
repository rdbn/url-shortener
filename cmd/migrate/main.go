package main

import (
	"flag"
	"log"
	"url-shortener/config"
	"url-shortener/pkg/migrate"
)

func main() {
	var (
		command      string
		forceVersion int
	)

	flag.StringVar(&command, "command", "up", "Migration command: up, down, force")
	flag.IntVar(&forceVersion, "version", 1, "Version to force (for force command)")
	flag.Parse()

	cfg := config.Load()

	switch command {
	case "up":
		log.Println("Running migrations up...")

		err := migrate.RunMigrations(cfg.GetPostgresDSN())
		if err != nil {
			log.Fatal("Migration failed:", err)
		}
	case "down":
		log.Println("Running migrations down...")
	case "force":
		log.Printf("Forcing migration version to %d...", forceVersion)

		err := migrate.ForceVersion(cfg.GetPostgresDSN(), forceVersion)
		if err != nil {
			log.Fatal("Force version failed:", err)
		}
	default:
		log.Fatal("Unknown command:", command)
	}
}
