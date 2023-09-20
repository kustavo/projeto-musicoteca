package internal

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
)

func Transfer(sourceFiles map[string]Artist, destinyFiles map[string]Artist, includeAudio bool, includeVideo bool, flags []string) error {
	numCPUs := runtime.NumCPU() // Obtém o número de CPUs disponíveis
	runtime.GOMAXPROCS(numCPUs) // Define o máximo de CPUs a serem utilizadas

	fmt.Println("Converting files to mp3...")

	var wg sync.WaitGroup
	for _, sourceArtist := range sourceFiles {
		for _, sourceAlbums := range sourceArtist.albums {
			for _, sourceSong := range sourceAlbums.songs {
				_ = sourceSong
			}
		}

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
	}

	wg.Wait()
	return nil
}

func convertFlacToMp3(flacPath string, mp3Path string) error {
	cmd := exec.Command("cp", flacPath, mp3Path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert FLAC to MP3: %v", err)
	}
	return nil
}

// func convertFlacToMp3(flacPath string, mp3Path string) error {
// 	cmd := exec.Command("ffmpeg", "-i", flacPath, "-y", "-ab", "320k", mp3Path)
// 	err := cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("failed to convert FLAC to MP3: %v", err)
// 	}
// 	return nil
// }
