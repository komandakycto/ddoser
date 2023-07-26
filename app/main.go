package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
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

	// create ticker and read the file every opts.ReadIntervalSeconds seconds
	t := time.NewTicker(time.Duration(opts.ReadIntervalSeconds) * time.Second)

doneApp:
	for {
		select {
		case <-t.C:
			l.Infof("Reading the file %s", opts.LogPath)
			// read the last N bytes from the file
			lines, err := readEndBytes(opts.LogPath, opts.BytesToRead, l)
			if err != nil {
				l.WithError(err).Errorf("Error reading file %s", opts.LogPath)
				continue
			}
			if len(lines) == 0 {
				l.Info("The log file is empty.")
				continue
			}

			// will process the lines concurrently, for this we need to split the lines into groups
			groups := splitSliceIntoGroups(lines, opts.LinesInGroup)
			// process the groups concurrently
			result := processGroupsConcurrently(ctx, groups, opts.IpNumbersThreshold, opts.TimeWindow, opts.UrlPattern, opts.JsonLogFormat, opts.OnlyIpv4, l)
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

	if overwrite {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	} else {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}

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
