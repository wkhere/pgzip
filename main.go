// pgzip decompresses or compresses gzip stream in parallel mode.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/klauspost/pgzip"
)

type config struct {
	compress      bool
	compressLevel int

	help func()
}

func run(conf config) error {
	switch {
	case conf.compress:
		w, err := pgzip.NewWriterLevel(os.Stdout, conf.compressLevel)
		if err != nil {
			return fmt.Errorf("failed creating pgzip writer: %w", err)
		}
		_, err = io.Copy(w, os.Stdin)
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("compress closing: %w", err)
		}

	case !conf.compress:
		r, err := pgzip.NewReader(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed creating pgzip reader: %w", err)
		}
		defer r.Close()
		_, err = io.Copy(os.Stdout, r)
		if err != nil {
			return fmt.Errorf("decompress: %w", err)
		}
	}

	return nil
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if conf.help != nil {
		conf.help()
		os.Exit(0)
	}

	err = run(conf)
	if err != nil {
		die(1, err)
	}
}

func die(code int, err error) {
	fmt.Fprintln(os.Stderr, "pgzip:", err)
	os.Exit(code)
}
