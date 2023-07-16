package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ArtistsFolder = "artistsFolder"
	ArtistFolder  = "artistFolder"
	AlbumFolder   = "albumFolder"
)

const (
	AudioMediaType = "audio"
	VideoMediaType = "video"
)

var IgnoredArtistsFolderFiles = []string{"Readme.md", ".sync.ffs_db"}

type song struct {
	name      string
	track     int
	flags     map[string]bool
	path      string
	midiaType string
}

type album struct {
	name  string
	flags map[string]bool
	songs map[string]song
	path  string
}

type artist struct {
	name   string
	albums map[string]album
	path   string
}

func main() {
	rootDir := flag.String("root-dir", "/mnt/arquivos/Músicas", "Root directory path")
	sourceDir := flag.String("source-dir", "/mnt/arquivos/Músicas", "Source directory path")
	outputDir := flag.String("output-dir", "", "Output directory path")
	includeAudio := flag.Bool("include-audio-files", false, "Include audio files")
	includeVideo := flag.Bool("include-video-files", false, "Include video files")

	flag.Parse()

	if sourceDir == nil || *sourceDir == "" {
		log.Fatal("source directory path not specified")
		return
	}

	if outputDir == nil || *outputDir == "" {
		log.Fatal("output directory path not specified")
		return
	}

	if rootDir == nil || *rootDir == "" {
		rootDir = sourceDir
		return
	}

	if !*includeAudio && !*includeVideo {
		*includeAudio = true
	}

	files, err := listFilesAndCheckEmptyDirs(*rootDir, *sourceDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file)
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

func listFilesAndCheckEmptyDirs(rootDir string, sourceDir string) (map[string]artist, error) {
	songs := make(map[string]artist)

	err := filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			isEmpty, err := isDirectoryEmpty(path)
			if err != nil {
				return err
			}
			if isEmpty {
				return fmt.Errorf("empty directory found: %s", path)
			}
		} else {
			songObj := song{}
			albumObj := album{}
			artistObj := artist{}

			relativePath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}

			arrRelativePath := strings.Split(relativePath, string(filepath.Separator))

			if arrRelativePath[0] == "_000_analise" || arrRelativePath[0] == "_000_fila" {
				return nil
			}

			if len(arrRelativePath) == 1 { // artists folder
				if !contains(IgnoredArtistsFolderFiles, arrRelativePath[0]) {
					return fmt.Errorf("invalid file location: %s", relativePath)
				}
				return nil
			} else if len(arrRelativePath) == 2 { // artist folder
				artistObj.name = arrRelativePath[0]
				artistObj.path = filepath.Dir(path)
				albumObj.name = ""
				albumObj.path = ""
				songObj.name = arrRelativePath[1]
			} else if len(arrRelativePath) == 3 { // album folder
				artistObj.name = arrRelativePath[0]
				artistObj.path = filepath.Dir(filepath.Dir(path))
				albumObj.name = arrRelativePath[1]
				albumObj.path = filepath.Dir(path)
				songObj.name = arrRelativePath[2]
			} else {
				return fmt.Errorf("directory beyond depth limit: %s", path)
			}
			appendMap(songs, artistObj, albumObj, songObj)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func appendMap(m map[string]artist, newArtist artist, newAlbum album, newSong song) {
	artistObj, artistExists := m[newArtist.name]
	if !artistExists {
		newArtist.albums = make(map[string]album)
		artistObj = newArtist
		m[newArtist.name] = artistObj
	}

	albumMap := artistObj.albums
	albumObj, albumExists := albumMap[newAlbum.name]
	if !albumExists {
		newAlbum.songs = make(map[string]song)
		albumObj = newAlbum
		albumMap[newAlbum.name] = albumObj
	}

	songMap := albumObj.songs
	if _, songExists := songMap[newSong.name]; !songExists {
		songMap[newSong.name] = newSong
	}
}

func isDirectoryEmpty(path string) (bool, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			isEmpty, err := isDirectoryEmpty(filepath.Join(path, entry.Name()))
			if err != nil {
				return false, err
			}
			if !isEmpty {
				return false, nil
			}
		} else {
			return false, nil
		}
	}

	return true, nil
}

func contains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
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

// func convertFlacToMp3(flacPath string, mp3Path string) error {

// 	// cmd := exec.Command("ffmpeg", "-i", flacPath, "-y", mp3Path)
// 	cmd := exec.Command("ffmpeg", "-i", flacPath, "-y", "-ab", "320k", mp3Path)
// 	err := cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("failed to convert FLAC to MP3: %v", err)
// 	}
// 	return nil
// }

// func getFolderLevelType(searchType string, relativePath string) (string, error) {
// 	numSeparators := len(filepath.SplitList(relativePath)) - 1

// 	switch searchType {
// 	case ArtistsSearchType:
// 		switch numSeparators {
// 		case 0:
// 			return ArtistsSearchType, nil
// 		case 1:
// 			return ArtistSearchType, nil
// 		case 2:
// 			return AlbumSearchType, nil
// 		default:
// 			return "", fmt.Errorf("directory beyond depth limit: %s", relativePath)
// 		}
// 	case ArtistSearchType:
// 		switch numSeparators {
// 		case 0:
// 			return ArtistSearchType, nil
// 		case 1:
// 			return AlbumSearchType, nil
// 		default:
// 			return "", fmt.Errorf("directory beyond depth limit: %s", relativePath)
// 		}
// 	case AlbumSearchType:
// 		switch numSeparators {
// 		case 0:
// 			return AlbumSearchType, nil
// 		default:
// 			return "", fmt.Errorf("directory beyond depth limit: %s", relativePath)
// 		}
// 	default:
// 		return "", fmt.Errorf("unknown searchType: %s", searchType)
// 	}
// }
