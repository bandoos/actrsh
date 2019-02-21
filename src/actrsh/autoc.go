package main

import (
	"bufio"
	"github.com/chzyer/readline"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// === Completer Builder ===

// NewComplt builds the readline.AutoCompleter instance for this shell
// TODO add custom config for static additions
// accepts a list of paths to get autocompletion for function defintions `paths`
// and a directory path to get autocompletion for the `(load "<path>") instruction
func NewComplt(paths []string, modelsDirs []string) readline.AutoCompleter {
	var completer *readline.PrefixCompleter
	completer = readline.NewPrefixCompleter(
		// the `:!` prefix will be used in further releases for extra commands
		readline.PcItem(":!",
			readline.PcItem("help"),
			readline.PcItem("show"),
		),
		// the `(load` prefix trigges dynamic autocompletion
		// for filenames in modelsDir
		readline.PcItem("(load",
			readline.PcItemDynamic(listFiles(modelsDirs)),
		),
		// the `(` prefix triggers dynamic autocompletion
		// for function names in the current selected model
		readline.PcItem("(",
			readline.PcItemDynamic(makeFuncsCompleter(paths)),
			readline.PcItem("run"),
		),
	)
	return completer
}

// ===== Model path completion logic ======

// listFiles return a closure that lists (in a lazy eval fashion)
// the files in the given path. This is intended to be the
// autocomplete source for the `(load` prefix
// NOTE this function adds a leading `"` to file names
// for convenience in the actrsh
func listFiles(paths []string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		for _, path := range paths {
			files, _ := ioutil.ReadDir(path)
			for _, f := range files {
				abs, err := filepath.Abs(path)
				panicErr(err)
				if filepath.Ext(f.Name()) == ".lisp" {
					names = append(names, "\""+filepath.Join(abs, f.Name())+"\"")
				}
			}
		}
		return names
	}
}

// ==== Function name completion logic ====

// makeFuncsCompleter returns the closure to be used as
// completer function for the `(` prefix. Invokes getFuncs(path)
// for each path in paths and joins all results into "names"
func makeFuncsCompleter(paths []string) func(string) []string {
	names := make([]string, 0)
	for _, path := range paths {
		names = append(names, getFuncs(path)...)
	}

	return func(line string) []string {
		return names
	}
}

// extractFunc is a helper used by getFuncs
// to parse a line and extract from it the first token following the
// first occurence of the `key` token passed as arg
// NOTE that it uses strings.toLower to render the process case insesitive
func extractFunc(line []byte, key string) string {
	tokens := strings.Split(string(line), " ")
	for i, t := range tokens {
		t = strings.ToLower(t)
		if strings.Contains(t, key) && i+1 < len(tokens) {
			return tokens[i+1]
		}
	}
	return ""
}

// getFuncs extract the string tokens which are destined
// to be autocompletion elements of the shell eventually.
func getFuncs(path string) []string {
	var file *os.File
	file, err := os.Open(path)
	panicErr(err)
	var funcs []string
	funcs = make([]string, 0)
	rdr := bufio.NewReader(file)
	for {
		line, err := rdr.ReadBytes('\n')
		if err == io.EOF {
			return funcs
		}
		panicErr(err)
		funcName := extractFunc(line, "defun")
		if funcName != "" {
			funcs = append(funcs, funcName)
		}
	}
}
