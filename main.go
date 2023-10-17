package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kustavo/projeto-musicoteca/internal"
)

func main() {
	argOrigem := flag.String("origem", "/mnt/arquivos/Músicas", "Diretório dos artistas na origem")
	argDestino := flag.String("destino", "", "Diretório dos artistas no destino")
	argArtista := flag.String("artista", "", "Especificar artista")
	argSomenteAudio := flag.Bool("somente-audio", false, "Somente arquivos de áudio")
	argSomenteVideo := flag.Bool("somente-video", false, "Somente arquivos de vídeo")
	argFlagsFiltro := flag.String("flags", "", "Flags")

	flag.Parse()

	if argDestino == nil || *argDestino == "" {
		log.Fatal("local do destino não informado")
		return
	}

	fmt.Println("Mapeando origem...")
	mapOrigem, err := internal.Mapear(*argOrigem, *argArtista)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mapeando destino...")
	mapDestino, err := internal.Mapear(*argDestino, *argArtista)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Analisando remoção arquivos do destino...")
	itens := internal.ObterArquivosRemocao(mapDestino, mapOrigem)
	if len(itens) > 0 {
		for _, item := range itens {
			fmt.Println("    - Remover: " + item)
		}
		fmt.Println("Deseja remover os arquivos acima do destino? (y/n)")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if input == "y\n" {
			fmt.Println("Removendo arquivos no destino")
			for _, item := range itens {
				err := os.RemoveAll(item)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	fmt.Println("Analisando envio de arquivos para o destino...")
	err = internal.Filtrar(mapOrigem, *argSomenteAudio, *argSomenteVideo, *argFlagsFiltro)
	if err != nil {
		log.Fatal(err)
	}
	itens = internal.ObterArquivosInclusao(mapOrigem, mapDestino)

	fmt.Println("Transferindo arquivos para o destino...")
	internal.Transfer(itens, *argOrigem, *argDestino)

	fmt.Println("Fim!")
}
