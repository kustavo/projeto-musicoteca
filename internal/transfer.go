package internal

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func Transfer(sourceFiles map[string]Artist, destinyPath string, includeAudio bool, includeVideo bool, flags []string) error {
	numCPUs := runtime.NumCPU() // Obtém o número de CPUs disponíveis
	runtime.GOMAXPROCS(numCPUs) // Define o máximo de CPUs a serem utilizadas

	fmt.Println("Converting files to mp3...")

	var wg sync.WaitGroup
	for _, sourceArtist := range sourceFiles {
		for _, sourceAlbums := range sourceArtist.albums {
			for _, sourceSong := range sourceAlbums.songs {

				// Obtém o diretório pai
				basePath := filepath.Dir(sourceArtist.path)
				filePath := filepath.Join(basePath, sourceSong.path)
				file := filepath.Base(filePath)

				if sourceSong.midiaType == AudioMediaType {
					mp3Path := filepath.Join(destinyPath, strings.TrimSuffix(sourceSong.path, ".flac")+".mp3")

					wg.Add(1)
					go func(filePath, mp3Path string) {
						defer wg.Done()
						fmt.Println(filePath)
						err := convertFlacToMp3(filePath, mp3Path)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Printf("Converted %s to %s", filePath, mp3Path)
						}
					}(filePath, mp3Path)

					fmt.Println(file)
				} else {
					videoPath := filepath.Join(destinyPath, sourceSong.path)
					moveFile(filePath, videoPath)
				}
			}
		}
	}

	wg.Wait()
	return nil
}

func moveFile(sourceFilePath string, destFilePath string) error {
	cmd := exec.Command("mkdir", "-p", filepath.Dir(destFilePath))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create dir: %v", err)
	}

	cmd = exec.Command("cp", sourceFilePath, destFilePath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to copy: %v", err)
	}
	return nil
}

func convertFlacToMp3(flacPath string, mp3Path string) error {
	cmd := exec.Command("mkdir", "-p", filepath.Dir(mp3Path))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create dir: %v", err)
	}

	cmd = exec.Command("ffmpeg", "-i", flacPath, "-y", "-ab", "320k", mp3Path)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert FLAC to MP3: %v", err)
	}
	return nil
}
