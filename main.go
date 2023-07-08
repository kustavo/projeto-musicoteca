package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func main() {

	dirIn := flag.String("d", "/mnt/arquivos/Músicas", "Input Directory")
	dirOut := flag.String("o", "", "Output Directory")
	typeAudio := flag.Bool("a", false, "Audio files")
	typeVideo := flag.Bool("v", false, "Vídeo files")

	flag.Parse()

	if dirOut == nil || *dirOut == "" {
		log.Fatal("output directory not found")
		return
	}

	if !*typeAudio && !*typeVideo {
		*typeAudio = true
	}

	files, err := getFilesWithSubstring(*dirIn, "[TOP]")
	if err != nil {
		log.Fatal(err)
	}

	numCPUs := runtime.NumCPU() // Obtém o número de CPUs disponíveis
	runtime.GOMAXPROCS(numCPUs) // Define o máximo de CPUs a serem utilizadas
	var wg sync.WaitGroup

	fmt.Println("Converting files to mp3...")
	for _, file := range files {
		mp3Path := filepath.Join(*dirOut, strings.TrimSuffix(file, ".flac")+".mp3")
		flacPath := filepath.Join(*dirIn, file)

		wg.Add(1)
		go func(flacPath, mp3Path string) {
			defer wg.Done()
			fmt.Println(flacPath)
			err = convertFlacToMp3(flacPath, mp3Path)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Converted %s to %s", flacPath, mp3Path)
			}
		}(flacPath, mp3Path)

		fmt.Println(file)
	}

	wg.Wait()
}

func getFilesWithSubstring(path string, substring string) ([]string, error) {
	var filesWithSubstring []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), substring) && strings.ToLower(filepath.Ext(info.Name())) == ".flac" {
			filesWithSubstring = append(filesWithSubstring, info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesWithSubstring, nil
}

func convertFlacToMp3(flacPath string, mp3Path string) error {

	// cmd := exec.Command("ffmpeg", "-i", flacPath, "-y", mp3Path)
	cmd := exec.Command("ffmpeg", "-i", flacPath, "-y", "-ab", "320k", mp3Path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert FLAC to MP3: %v", err)
	}
	return nil
}
