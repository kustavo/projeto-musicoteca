package main

import (
	"flag"
	"log"
	"strings"

	"github.com/kustavo/projeto-musicoteca/internal"
)

func main() {
	rootPath := flag.String("root-path", "/mnt/arquivos/Músicas", "Root directory path")
	sourcePath := flag.String("source-path", "/mnt/arquivos/Músicas", "Source directory path")
	outputPath := flag.String("output-path", "", "Output directory path")
	includeAudio := flag.Bool("include-audio-files", false, "Include audio files")
	includeVideo := flag.Bool("include-video-files", false, "Include video files")
	flagsFilter := flag.String("flags", "*", "Flags")

	flag.Parse()

	if sourcePath == nil || *sourcePath == "" {
		log.Fatal("source directory path not specified")
		return
	}

	if outputPath == nil || *outputPath == "" {
		log.Fatal("output directory path not specified")
		return
	}

	if rootPath == nil || *rootPath == "" {
		rootPath = sourcePath
		return
	}

	if !*includeAudio && !*includeVideo {
		*includeAudio = true
	}

	sourceFiles, err := internal.GetFilesMap(*rootPath, *sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	destinyFiles, err := internal.GetFilesMap(*outputPath, *outputPath)
	if err != nil {
		log.Fatal(err)
	}

	flags := strings.Split(*flagsFilter, ",")
	internal.Filter(sourceFiles, destinyFiles, *includeAudio, *includeVideo, flags)
	internal.Transfer(sourceFiles, *outputPath, *includeAudio, *includeVideo, flags)
}
