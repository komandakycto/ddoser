package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"

	"github.com/komandakycto/ddoser/ddoser/helpers"
	"github.com/komandakycto/ddoser/ddoser/parser/nginx"
	"github.com/komandakycto/ddoser/ddoser/processor"
	"github.com/komandakycto/ddoser/ddoser/reader"
)

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.Parse(); err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	// init logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	l := logrus.NewEntry(log)

	// Capture system signals.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// app context
	ctx := context.Background()

	l.Info("Starting the application...")

	// create log file logReader.
	logReader := reader.NewLogReader(opts.LogPath, opts.BytesToRead, opts.AverageLogBytes)
	// create log row parser.
	jsonTimeLayout := time.RFC3339
	if opts.JsonLogTimeLayout != "" {
		jsonTimeLayout = opts.JsonLogTimeLayout
	}
	parser := nginx.NewParser(!opts.JsonLogFormat, jsonTimeLayout)
	// create rows ipProcessor.
	ipProcessor := processor.NewIPAnalysis(opts.IpNumbersThreshold, opts.TimeWindow, opts.UrlPattern, parser, opts.OnlyIpv4, l)

	// create ticker and read the file every opts.ReadIntervalSeconds seconds
	t := time.NewTicker(time.Duration(opts.ReadIntervalSeconds) * time.Second)

doneApp:
	for {
		select {
		case <-t.C:
			l.Infof("Reading the file %s", opts.LogPath)
			// read the last N bytes from the file
			lines, err := logReader.ReadLastBytes()
			if err != nil {
				l.WithError(err).Errorf("Error reading file %s", opts.LogPath)
				continue
			}
			if len(lines) == 0 {
				l.Info("The log file is empty.")
				continue
			}

			// will process the lines concurrently, for this we need to split the lines into groups
			groups := helpers.SplitSliceIntoGroups(lines, opts.LinesInGroup)
			// process the groups concurrently
			result, err := ipProcessor.Process(ctx, groups)
			if err != nil {
				l.WithError(err).Errorf("Error processing the file %s", opts.LogPath)
				continue
			}
			// write the result to the output file
			err = writeResult(result, opts.OutputPath, opts.OutputOverwrite)
			if err != nil {
				l.WithError(err).Errorf("Error writing to file %s", opts.OutputPath)
				continue
			}
		case <-signalCh:
			t.Stop()
			ctx.Done()
			l.Info("Received a signal to stop.")
			break doneApp
		}
	}
}

func writeResult(data map[string]bool, filePath string, overwrite bool) error {
	// Open the file with write-only and create/truncate options (to overwrite existing file or create a new one).
	var file *os.File
	var err error

	mode := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if overwrite {
		mode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	file, err = os.OpenFile(filePath, mode, 0644)

	if err != nil {
		return fmt.Errorf("failed to open the file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Write each element of the slice to the file (one element per line).
	for ip, _ := range data {
		_, err := fmt.Fprintln(file, ip)
		if err != nil {
			return fmt.Errorf("failed to write to the file: %v", err)
		}
	}

	return nil
}
