package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/sourcegraph/conc/stream"
	"github.com/spf13/pflag"
	"go.chensl.me/redglob"
)

var (
	stdout = bufio.NewWriter(os.Stdout)

	withFilename = pflag.BoolP(
		"with-filename",
		"H",
		false,
		"print file name with output lines",
	)
	lineNumber = pflag.BoolP(
		"line-number",
		"n",
		false,
		"print line number with output lines",
	)
	invertMatch = pflag.BoolP("invert-match", "v", false, "select non-matching lines")
	maxProcs    = pflag.IntP(
		"max-procs",
		"P",
		runtime.GOMAXPROCS(-1),
		"run at most MAX-PROCS processes at a time",
	)
	help = pflag.BoolP("help", "h", false, "display this help and exit")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Search for PATTERNS in each FILE\n\n")
	fmt.Fprintf(os.Stderr, "Usage: grep [OPTION]... PATTERNS [FILE]...\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	pflag.PrintDefaults()
}

func main() {
	log.SetPrefix("ggrep: ")
	log.SetFlags(0)
	pflag.Usage = usage
	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}

	nArg := pflag.NArg()
	if nArg < 1 {
		pflag.Usage()
		os.Exit(2)
	}

	defer stdout.Flush()

	pat := pflag.Arg(0)
	files := pflag.Args()[1:]

	if len(files) == 0 {
		if err := grep(pat, "(standard input)", os.Stdin); err != nil {
			log.Fatal(err)
		}
		return
	}
	if len(files) > 1 {
		*withFilename = true
	}

	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		if err := grep(pat, filename, f); err != nil {
			log.Fatal(err)
		}
		_ = f.Close()
	}
}

func grep(pat, filename string, f io.Reader) error {
	s := stream.New().WithMaxGoroutines(*maxProcs)

	r := bufio.NewReader(f)
	var i int64
	for {
		l, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if l[len(l)-1] == '\n' {
			l = l[:len(l)-1]
		}
		i++
		i := i
		s.Go(func() stream.Callback {
			matched := redglob.MatchBytes(l, pat)
			if *invertMatch {
				matched = !matched
			}
			return func() {
				if !matched {
					return
				}
				if *withFilename {
					stdout.WriteString(filename)
					stdout.WriteByte(':')
				}
				if *lineNumber {
					stdout.WriteString(strconv.FormatInt(i, 10))
					stdout.WriteByte(':')
				}
				stdout.Write(l)
				stdout.WriteByte('\n')
			}
		})
	}

	s.Wait()
	return nil
}
