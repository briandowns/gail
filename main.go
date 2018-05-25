// build +freebsd
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
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
	pathFlag    string
)

const defaultJailPath = "/zroot/jails"

const usage = `version: %s - git: %s

Usage: %s [-v] [-h] [-i]

Options:
  -h            help
  -v            version
  -p            path to jail fs

Examples: 
  %[3]s ls                     List contents of current directory
`

func usageFunc() {
	w := os.Stderr
	for _, arg := range os.Args {
		if arg == "-h" {
			w = os.Stdout
			break
		}
	}
	fmt.Fprintf(w, usage, version, gitSHA, name)
}

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		for range sc {
			os.Exit(1)
		}
	}()
	if len(os.Args) < 2 {
		usageFunc()
		os.Exit(1)
	}

	flag.Usage = usageFunc
	flag.BoolVar(&versionFlag, "v", false, "")
	flag.StringVar(&pathFlag, "p", defaultJailPath, "")
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
		Path:     path + "/build",
		Hostname: id,
		Name:     id,
		Chdir:    true,
	}
	if _, err := jail.Jail(&opts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var cmd *exec.Cmd
	if len(os.Args) > 2 {
		cmd = exec.Command(os.Args[1], os.Args[2:]...)
	} else {
		cmd = exec.Command(os.Args[1])
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
