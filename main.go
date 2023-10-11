package main

import (
	"flag"
	"log"

	"github.com/kustavo/projeto-musicoteca/internal"
)

func main() {
	rootPath := flag.String("root-path", "/mnt/arquivos/Músicas", "Root directory path")
	sourcePath := flag.String("source-path", "/mnt/arquivos/Músicas", "Source directory path")
	outputPath := flag.String("output-path", "", "Output directory path")
	includeAudio := flag.Bool("include-audio-files", false, "Include audio files")
	includeVideo := flag.Bool("include-video-files", false, "Include video files")

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

	// TODO: Filtrar por tags

	flags := []string{"ROC"}
	internal.Filter(sourceFiles, destinyFiles, *includeAudio, *includeVideo, flags)
	internal.Transfer(sourceFiles, *outputPath, *includeAudio, *includeVideo, flags)

	for _, file := range sourceFiles {
		// fmt.Println(file)
		_ = file
	}

	// files, err := getFilesWithSubstring(*dirIn, "[TOP]")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// numCPUs := runtime.NumCPU() // Obtém o número de CPUs disponíveis
	// runtime.GOMAXPROCS(numCPUs) // Define o máximo de CPUs a serem utilizadas
	// var wg sync.WaitGroup

	// fmt.Println("Converting files to mp3...")
	// for _, file := range files {
	// 	mp3Path := filepath.Join(*dirOut, strings.TrimSuffix(file, ".flac")+".mp3")
	// 	flacPath := filepath.Join(*dirIn, file)

	// 	wg.Add(1)
	// 	go func(flacPath, mp3Path string) {
	// 		defer wg.Done()
	// 		fmt.Println(flacPath)
	// 		err = convertFlacToMp3(flacPath, mp3Path)
	// 		if err != nil {
	// 			log.Println(err)
	// 		} else {
	// 			log.Printf("Converted %s to %s", flacPath, mp3Path)
	// 		}
	// 	}(flacPath, mp3Path)

	// 	fmt.Println(file)
	// }

	// wg.Wait()
}

// func getFilesWithSubstring(path string, substring string) ([]string, error) {
// 	var filesWithSubstring []string

// 	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if !info.IsDir() && strings.Contains(info.Name(), substring) && strings.ToLower(filepath.Ext(info.Name())) == ".flac" {
// 			filesWithSubstring = append(filesWithSubstring, info.Name())
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return filesWithSubstring, nil
// }
