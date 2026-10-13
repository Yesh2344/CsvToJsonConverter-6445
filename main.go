package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/example/csvtojsonconverter/config"
	"github.com/example/csvtojsonconverter/converter"
	"github.com/example/csvtojsonconverter/handlers"
	"github.com/example/csvtojsonconverter/logger"
	"go.uber.org/zap"
)

func main() {
	// Command‑line flags
	var (
		csvPath  string
		jsonPath string
		serve    bool
		port     int
	)
	flag.StringVar(&csvPath, "input", "", "Path to input CSV file")
	flag.StringVar(&jsonPath, "output", "", "Path to output JSON file")
	flag.BoolVar(&serve, "serve", false, "Start HTTP server")
	flag.IntVar(&port, "port", 8080, "Port for HTTP server")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	// Initialise logger
	if err := logger.Init(cfg.LogLevel); err != nil {
		fmt.Fprintf(os.Stderr, "logger init error: %v\n", err)
		os.Exit(1)
	}
	defer logger.Logger.Sync()

	if serve {
		http.HandleFunc("/convert", handlers.ConvertHandler(cfg))
		addr := fmt.Sprintf(":%d", port)
		logger.Logger.Info("starting HTTP server", zap.String("addr", addr))
		if err := http.ListenAndServe(addr, nil); err != nil {
			logger.Logger.Fatal("server failed", zap.Error(err))
		}
		return
	}

	// CLI mode – require both input and output
	if csvPath == "" || jsonPath == "" {
		fmt.Fprintln(os.Stderr, "both -input and -output flags are required in CLI mode")
		flag.Usage()
		os.Exit(1)
	}
	if err := converter.ConvertFile(csvPath, jsonPath, cfg); err != nil {
		logger.Logger.Error("conversion failed", zap.Error(err))
		os.Exit(1)
	}
}