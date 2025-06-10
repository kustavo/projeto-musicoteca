package sync

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/kustavo/projeto-musicoteca/internal/shared"
)

func Sync(medias []shared.MediaDTO, destination string, extDestination string, artistFolder bool) error {
	createDirectory(destination)
	songsToTransfer, filesToDelete := getFilesForSync(medias, destination, artistFolder, extDestination)

	log.Printf("Arquivos já sincronizados: %d", len(medias)-len(songsToTransfer))
	log.Printf("Arquivos a serem transferidos: %d", len(songsToTransfer))
	log.Printf("Arquivos a serem deletados: %d", len(filesToDelete))

	delete(filesToDelete)
	transfer(songsToTransfer, destination, extDestination, artistFolder)

	return nil
}

func getFilesForSync(medias []shared.MediaDTO, rootDestination string, artistFolder bool, extDestination string) ([]shared.MediaDTO, []shared.File) {
	var mediasToTransfer []shared.MediaDTO
	var filesToDelete []shared.File

	artistMediasMap := make(map[string][]shared.MediaDTO)
	for _, media := range medias {
		artistName := ""
		if artistFolder {
			artistName = media.ArtistName
		}
		artistMediasMap[artistName] = append(artistMediasMap[artistName], media)
	}

	for artistName, artistMedias := range artistMediasMap {
		dirDestination := rootDestination
		if artistFolder {
			dirDestination = filepath.Join(rootDestination, artistName)
		}
		filesDestination, err := shared.ListDirectoryFiles(dirDestination)
		if err != nil {
			log.Printf("Erro ao listar arquivos no destino: %v", err)
			continue
		}
		filesDestination = shared.FilterFilesByExtension(filesDestination, extDestination)

		fileSet := make(map[string]shared.File)
		for _, file := range filesDestination {
			fileSet[file.Name] = file
		}
		for _, media := range artistMedias {
			destFile, found := fileSet[media.Name]
			if !found {
				mediasToTransfer = append(mediasToTransfer, media)
				continue
			}
			c, err := compareFileTimestamps(media.Path, destFile.Path)
			if err != nil {
				log.Printf("Erro ao comparar timestamps: %v", err)
				continue
			}
			if c != 0 {
				mediasToTransfer = append(mediasToTransfer, media)
			}
		}

		for _, file := range filesDestination {
			found := false
			for _, song := range artistMedias {
				if file.Name == song.Name {
					found = true
					break
				}
			}
			if !found {
				filesToDelete = append(filesToDelete, file)
			}
		}
	}

	return mediasToTransfer, filesToDelete
}

func delete(files []shared.File) {
	for _, file := range files {
		log.Printf("Deletando: %s", file.Path)
		err := deleteFile(file.Path)
		if err != nil {
			log.Printf("Erro ao deletar %s: %v", file.Path, err)
		}
	}
}

func transfer(medias []shared.MediaDTO, rootDestination string, extDestination string, artistFolder bool) error {
	numWorkers := shared.Conf.NumWorkers
	sem := make(chan struct{}, numWorkers)
	var wg sync.WaitGroup

	for _, media := range medias {
		dirDestination := rootDestination
		if artistFolder {
			dirDestination = filepath.Join(rootDestination, media.ArtistName)
		}

		createDirectory(dirDestination)
		filePathSource := media.Path
		filePathDestination := filepath.Join(dirDestination, media.Name+extDestination)

		switch extDestination {
		case ".mp3":
			if media.Extension == ".flac" {
				wg.Add(1)
				sem <- struct{}{}
				go func(filePathSource, filePathDestination string) {
					defer wg.Done()
					defer func() { <-sem }()
					log.Printf("Convertendo: %s", filePathSource)
					err := convertFlacToMp3(filePathSource, filePathDestination)
					if err != nil {
						log.Printf("Erro ao converter %s: %v", filePathSource, err)
					}
				}(filePathSource, filePathDestination)
			} else {
				log.Printf("Não há conversor para a extensão: %s", media.Extension)
			}
		default:
			if media.Extension == extDestination {
				log.Printf("Copiando: %s", filePathSource)
				err := copyFile(filePathSource, filePathDestination)
				if err != nil {
					return err
				}
			} else {
				log.Printf("Não há conversor para %s -> %s", media.Extension, extDestination)
			}
		}
	}

	wg.Wait()
	return nil
}

func compareFileTimestamps(path1 string, path2 string) (int, error) {
	info1, err := os.Stat(path1)
	if err != nil {
		return 0, fmt.Errorf("erro ao obter informações do arquivo: %v", err)
	}
	info2, err := os.Stat(path2)
	if err != nil {
		return 0, fmt.Errorf("erro ao obter informações do arquivo: %v", err)
	}

	return info1.ModTime().Compare(info2.ModTime()), nil
}

func deleteFile(path string) error {
	cmd := exec.Command("rm", path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao deletar arquivo: %v", err)
	}
	return nil
}

func createDirectory(path string) error {
	cmd := exec.Command("mkdir", "-p", path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao criar diretório: %v", err)
	}
	return nil
}

func copyFile(sourcePath string, destinationPath string) error {
	cmd := exec.Command("cp", sourcePath, destinationPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao copiar: %v", err)
	}
	err = changeDates(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	return nil
}

func convertFlacToMp3(sourcePath string, destinationPath string) error {
	cmd := exec.Command("ffmpeg", "-i", sourcePath, "-y", "-ab", "320k", destinationPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao converter FLAC para MP3: %v", err)
	}
	err = changeDates(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	return nil
}

func changeDates(sourcePath string, destinationPath string) error {
	cmd := exec.Command("touch", "-r", sourcePath, destinationPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao alterar datas do arquivo: %v", err)
	}
	return nil
}
