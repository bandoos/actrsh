// package main
// this is the executable cmd package for the actrsh
// This program is a wrapper for the original ACT-R standalone
// TODO put ref.
// The standalone offers an executable command that handles the
// startup bureacracy for ACT-R listener and environment
// Unfortunately the REPL offered by the default lisp coming
// with the standalone does not support lateral arrow navigation (to correct mistakes)
// nor vertical history scrolling (i.e. use arrows to navigate previous commands)
// The purpose of this package is enhanching the standalone with a richer shell
// that allows for the above mentioned features as well as static and dynamic completion.
// None of this would be possible without the great `readline` lib and
// the awesome io.Copy routine from golang stdlib
//
// NOTE that the program can actually be used to wrap any
// other program, the fact that i call it actrsh is because the autocomplete
// logic is tailored for act-r and particuarly its lisp version
//
// The core functionality here is patching stdout and stderr of another command (ref. 3rd party command)
// to this command's stdout,
// concurrently, we parse (with autocompletion and history features)
// line by line this cmd stdin and write the lines to the 3rd party command
package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/chzyer/readline"
	"io"
	"os"
	"os/exec"
)

const (
	BasePrompt string = "\033[31m»»\033[0m "
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var cli *readline.Instance
	var cmd *exec.Cmd

	// ==== FLAG PARSING ====
	parser := argparse.NewParser("actrsh_0.1", `(ACT-R shell)
A wrapper for the ACT-R standalone lisp REPL written in golang.

The wrapper adds:
	- history 
	- arrow navigation
	- autocompletion`)

	cmdpath := parser.String("c", "command", &argparse.Options{
		Required: true,
		Help:     "specify the 3rd party program to be executed in the background",
	})

	modelsDirs := parser.List("d", "models-dir", &argparse.Options{
		Required: false,
		Help:     `the path to the models directory  to AutoComplete '(load \"'. Can be repeated to specify more than one folder`,
	})

	models := parser.List("l", "model-list", &argparse.Options{
		Required: false,
		Help:     `can be repeated to specify more than one model. Functions parsed from the given paths will be added to the autocompletion for the "(" prefix`,
	})

	// parse flags
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	fmt.Println("models: ", *models)
	// recover internal error in case of panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[== Sorry there was a fatal error: ", r, "==]")
			os.Exit(1)
		}
	}()

	// ==== SETUP ====
	// build cli instance and cmd
	cli = NewCli(*models, *modelsDirs)
	cmd = exec.Command(*cmdpath)

	// setup stdin stdout and stder pipes
	fmt.Println("Setting up pipes")
	inp, err := cmd.StdinPipe()
	panicErr(err)

	out, err := cmd.StdoutPipe()
	panicErr(err)

	stderr, err := cmd.StderrPipe()
	panicErr(err)

	fmt.Println("Starting yor program")
	err = cmd.Start()
	panicErr(err)

	// ==== MAIN LOGIC ====
	// on a separate goroutine copy the sub-program's output to
	// this program's Stdout. io.Copy is extremely efficient :) gopher thanks stdlib!!
	go io.Copy(os.Stdout, out)
	go io.Copy(os.Stdout, stderr)
	// handle the input from this program's stdin in a third goroutine
	go func() {
		for {
			//fmt.Println("===> reading next command")
			bts, err := cli.ReadSlice()
			panicErr(err)
			inp.Write(bts)
		}
	}()

	// let the main goroutine wait for the target command end.
	err = cmd.Wait()
	panicErr(err)
	fmt.Println("[== The target program exited with status: ", err, "\n bye byte!! ==]")
}

// NewCli builds the cli instance for this shell
func NewCli(paths []string, modelsDirs []string) *readline.Instance {
	auto := NewComplt(paths, modelsDirs)
	l, err := readline.NewEx(&readline.Config{
		Prompt:            BasePrompt,
		HistoryFile:       "/tmp/actrsh.tmp",
		AutoComplete:      auto,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	panicErr(err)
	return l

}
