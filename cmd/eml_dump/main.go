package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fadend/eml"
)

func main() {
	inputFile := flag.String("input_eml", "", "Input .eml file")
	outputDir := flag.String("output_dir", "", "Output dir")
	flag.Parse()
	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "Missing required arg --input_eml")
		os.Exit(1)
	}
	if *outputDir == "" {
		fmt.Fprintln(os.Stderr, "Missing required arg --output_dir")
		os.Exit(1)
	}
	f, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	attachments, err := eml.ExtractFileNameToAttachment(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for name, content := range attachments {
		outputPath := filepath.Join(*outputDir, name)
		if err = os.WriteFile(outputPath, content, 0664); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
