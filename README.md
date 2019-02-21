# actrsh

### What 
actrsh is a small wrapper for the ACT-R standalone that aims at enhanching the user experience at the LISP shell for ACT-R.

### Why
The ACT-R standalone offers an executable command that handles the
startup bureacracy for ACT-R listener and environment. This is a bash script for \*nix systems.

Although this is really convenient and robust, unfortunately the shell (or rather REPL) offered by the default lisp coming
with the standalone does not support lateral arrow navigation (to correct mistakes)
nor vertical history scrolling (i.e. use arrows to navigate previous commands)

The purpose of this package is enhanching the standalone with a richer REPL
that allows for the above mentioned features as well as static and dynamic completion.

### How
None of this would be possible without the great [`readline`](https://github.com/chzyer/readline) lib and

the awesome Copy routine from golang stdlib [`io.Copy`](https://golang.org/src/io/io.go?s=12784:12844#L353)

The core functionality here is patching stdout and stderr of another command (ref. 3rd party command)
to this command's stdout. Concurrently, we parse (with autocompletion and history features)
line by line **this** cmd stdin and write the lines to the 3rd party command

## Features

By wrapping the ACT-R REPL with this utility we gain:

- history navigation via vertical arrow keys
- lateral movement with horizontal arrow keys

if the -l flag is specified at least once:

- experiment function autocompletion.
    - the program fill try to extract all (defun \<name\> ...) statements from the files specified by -l and build an autocompletion function for the `(` prefix that will scroll through available functions when tab is pressed.

if the -d flag is specified at least once:

- model laod autocompletion.
    - the program will iterate files in the directories specified by -d occurences (non recursively) and build an autocompletion function for the `(load` prefix that will scroll through available `*.lisp` files when tab is pressed.  

## Usage
```
usage: actrsh_0.1 [-h|--help] -c|--command "<value>" [-d|--models-dir "<value>"
                  [-d|--models-dir "<value>" ...]] [-l|--model-list "<value>"
                  [-l|--model-list "<value>" ...]]

                  (ACT-R shell)
A wrapper for the ACT-R standalone lisp REPL written in golang.

The wrapper adds:
	- history 
	- arrow navigation
	- autocompletion

Arguments:

  -h  --help        Print help information
  -c  --command     specify the 3rd party program to be executed in the
                    background
  -d  --models-dir  the path to the models directory  to AutoComplete the
                    '(load ' prefix. Can be repeated to specify more than one
                    folder
  -l  --model-list  can be repeated to specify more than one model. Functions
                    parsed from the given paths will be added to the
                    autocompletion for the "(" prefix

```

# Installing

if you are using an x86\_64 with linux kernel or x86\_64 mac then you can directly install the compiled bundles in this repo's `bin/` folder.

TODO instructions

If you are running on an ARM platform or Windows then you'll have to build from source manually

# Building 
In order to compile the project you must have the `go` command installed.
in your shell run `go version` to check that. If that fails the install go.

TODO instructions



# Examples

1. Start the act-r env. with the -c flag. Will decorate with history feature.
```
$ actrsh -c ~/ACT-R/run-act-r.command
```

2. Start env. and specify 2 files for function autocompletion
```
$ actrsh -c ~/path/to/ACT-R/run-act-r.command -l ~/path/to/sart.lisp -l ~/path/to/subitize.lisp
```

3. Start env. and specify 2 files for function autocompletion and directory for model loading autocompletion
```
$ actrsh -c ~/ACT-R/run-act-r.command -l ~/path/to/sart.lisp -l ~/pat/to/subitize.lisp -d ~/path/to/my-models/
```

