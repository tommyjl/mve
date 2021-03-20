package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var config struct {
	name    string
	editor  string
	verbose bool
}

func init() {
	config.name = "lsmv"
	config.editor = os.Getenv("EDITOR")
	if config.editor == "" {
		config.editor = "vim"
	}
	flag.BoolVar(&config.verbose, "v", false, "verbose")
}

func editor(filename string) *exec.Cmd {
	editor := exec.Command(config.editor, filename)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	return editor
}

func main() {
	flag.Parse()
	if config.verbose {
		log.SetFlags(log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	filenames := flag.Args()
	if len(filenames) < 1 {
		log.Fatal("No filenames provided")
	}

	tmpFile, err := ioutil.TempFile("", config.name)
	if err != nil {
		log.Fatal("Unable to create temp file")
	}
	defer os.Remove(tmpFile.Name())

	for _, filename := range filenames {
		tmpFile.WriteString(filename + "\n")
	}

	err = editor(tmpFile.Name()).Run()
	if err != nil {
		log.Fatal("Failed to open $EDITOR")
	}

	edited, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		log.Fatal("Unable to read temp file")
	}

	editedFilenames := strings.Split(string(edited), "\n")[:len(filenames)]
	if len(filenames) != len(editedFilenames) {
		log.Fatal("Deleting files is not supported")
	}

	for i, changed := range editedFilenames {
		original := filenames[i]
		if changed == original {
			continue
		}
		err = os.Rename(original, changed)
		if err != nil {
			log.Fatal("Unable to rename file")
		}

		log.Println(original, "->", changed)
	}
}
