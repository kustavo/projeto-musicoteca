package internal

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func Transfer(arquivos []string, origem string, destino string) error {
	numCPUs := runtime.NumCPU() // Obtém o número de CPUs disponíveis
	runtime.GOMAXPROCS(numCPUs) // Define o máximo de CPUs a serem utilizadas

	var wg sync.WaitGroup

	for _, arquivo := range arquivos {
		pathOrigem := arquivo
		pathDestino := strings.Replace(arquivo, origem, destino, 1)
		extensao := filepath.Ext(pathOrigem)

		switch extensao {
		case ".flac":
			pathDestino = strings.TrimSuffix(pathDestino, extensao) + obterExtensaoConvertida(extensao)
			fmt.Printf("Convertendo: %s >>> %s \n", pathOrigem, pathDestino)
			wg.Add(1)
			go func(pathOrigem, pathDestino string) {
				defer wg.Done()
				err := converterFlacParaMp3(pathOrigem, pathDestino)
				if err != nil {
					fmt.Println(err)
				}
			}(pathOrigem, pathDestino)
		default:
			fmt.Printf("Copiando: %s >>> %s \n", pathOrigem, pathDestino)
			err := copiarArquivo(pathOrigem, pathDestino)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	wg.Wait()
	return nil
}

func copiarArquivo(origem string, destino string) error {
	cmd := exec.Command("mkdir", "-p", filepath.Dir(destino))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao criar diretório: %v", err)
	}

	cmd = exec.Command("cp", origem, destino)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao copiar: %v", err)
	}
	return nil
}

func converterFlacParaMp3(origem string, destino string) error {
	cmd := exec.Command("mkdir", "-p", filepath.Dir(destino))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao criar diretório: %v", err)
	}

	cmd = exec.Command("ffmpeg", "-i", origem, "-y", "-ab", "320k", destino)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao converter FLAC para MP3: %v", err)
	}
	return nil
}
