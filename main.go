// pgzip decompresses or compresses gzip stream in parallel mode.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/klauspost/pgzip"
	"github.com/spf13/pflag"
)

const (
	prog      = "pgzip"
	usageHead = "Usage: " + prog + "\t(reads stdin, outputs to stdout)"
)

type config struct {
	decompress    bool
	compressLevel int
}

func parseFlags(args []string) config {
	var conf config
	var help bool

	flag := pflag.NewFlagSet("flags", pflag.ContinueOnError)
	flag.SortFlags = false

	flag.BoolVarP(&conf.decompress, "decompress", "d", false,
		"decompress (default false -- means compress)")
	flag.IntVarP(&conf.compressLevel, "level", "n", 6,
		"compress level, -2..9")
	flag.BoolVarP(&help, "help", "h", false,
		"show this help and exit")
	flag.Usage = func() {
		fmt.Fprintln(flag.Output(), usageHead)
		flag.PrintDefaults()
	}

	err := flag.Parse(args)
	if err != nil {
		flag.Usage()
		os.Exit(2)
	}
	if help {
		flag.SetOutput(os.Stdout)
		flag.Usage()
		os.Exit(0)
	}

	return conf
}

func main() {
	conf := parseFlags(os.Args[1:])

	switch {
	case !conf.decompress:
		w, err := pgzip.NewWriterLevel(os.Stdout, conf.compressLevel)
		if err != nil {
			die(fmt.Errorf("failed creating pgzip writer: %v", err))
		}
		_, err = io.Copy(w, os.Stdin)
		if err != nil {
			die(fmt.Errorf("compress: %v", err))
		}
		err = w.Close()
		if err != nil {
			die(fmt.Errorf("compress closing: %v", err))
		}

	default:
		r, err := pgzip.NewReader(os.Stdin)
		if err != nil {
			die(fmt.Errorf("failed creating pgzip reader: %v", err))
		}
		defer r.Close()
		_, err = io.Copy(os.Stdout, r)
		if err != nil {
			die(fmt.Errorf("decompress: %v", err))
		}
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, prog+":", err)
	os.Exit(1)
}
