package main

import (
	"flag"
	"log"

	"github.com/kustavo/projeto-musicoteca/internal"
)

func main() {
	argOrigem := flag.String("origem", "/mnt/arquivos/Músicas", "Diretório dos artistas na origem")
	argDestino := flag.String("destino", "", "Diretório dos artistas no destino")
	argArtista := flag.String("artista", "", "Especificar artista")
	// argSomenteAudio := flag.Bool("somente-audio", false, "Somente arquivos de áudio")
	// argSomenteVideo := flag.Bool("somente-video", false, "Somente arquivos de vídeo")
	// argFlagsFiltro := flag.String("flags", "", "Flags")

	flag.Parse()

	if argDestino == nil || *argDestino == "" {
		log.Fatal("local do destino não informado")
		return
	}

	artistas, err := internal.Mapear(*argOrigem, *argArtista)
	if err != nil {
		log.Fatal(err)
	}

	_ = artistas

	// destinyFiles, err := internal.GetFilesMap(*outputPath, *outputPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// flags := strings.Split(*flagsFilter, ",")
	// internal.Filter(sourceFiles, destinyFiles, *includeAudio, *includeVideo, flags)
	// internal.Transfer(sourceFiles, *outputPath, *includeAudio, *includeVideo, flags)
}
