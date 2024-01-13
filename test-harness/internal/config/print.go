package config

import (
	"fmt"
	"io"
)

func (cfg *TestsuiteEnvConfig) printLoadedFrom(w io.Writer) {
	cfgFiles := dotenvFiles(cfg.Env)
	fmt.Fprintln(w, "Loaded from:")
	for _, fileName := range cfgFiles {
		fmt.Fprintf(w, "\t- %s\n", fileName)
	}
	fmt.Fprintln(w, "")
}

func (cfg *TestsuiteEnvConfig) Print(w io.Writer) {
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "TEST SUITE CONFIG")
	fmt.Fprintln(w, "-----------------")

	cfg.printLoadedFrom(w)
}
