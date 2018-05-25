// build +freebsd
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/pborman/uuid"

	"github.com/briandowns/jail"
)

var (
	version string
	gitSHA  string
	name    string
)

var (
	versionFlag bool
	idFlag      bool
	pathFlag    string
)

const defaultJailPath = "/zroot/jails"

const usage = `version: %s - git: %s

Usage: %s [-v] [-h] [-i]

Options:
  -h            help
  -v            version
  -i            return JID
  -p            path to jail fs

Examples: 
  %[3]s ls                     List contents of current directory
  %[3]s -i ls                  List contents of current directory and return JID
`

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		for range sc {
			os.Exit(1)
		}
	}()

	flag.Usage = func() {
		w := os.Stderr
		for _, arg := range os.Args {
			if arg == "-h" {
				w = os.Stdout
				break
			}
		}
		fmt.Fprintf(w, usage, version, gitSHA, name)
	}
	flag.BoolVar(&versionFlag, "v", false, "")
	flag.StringVar(&pathFlag, "p", defaultJailPath, "")
	flag.BoolVar(&idFlag, "i", false, "")
	flag.Parse()
	if versionFlag {
		fmt.Fprintf(os.Stdout, "version: %s - %s\n", version, gitSHA)
		os.Exit(0)
	}
	id := uuid.NewUUID().String()
	var path string
	if pathFlag != "" {
		path = pathFlag
	} else {
		path = defaultJailPath
	}
	opts := jail.Opts{
		Chdir:    true,
		Hostname: id,
		Name:     id,
		Path:     path,
	}
	j, err := jail.Jail(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if idFlag {
		fmt.Println(j)
	}
	os.Exit(0)
}
