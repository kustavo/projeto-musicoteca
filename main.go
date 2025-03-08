package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kustavo/projeto-musicoteca/internal/mapping"
	"github.com/kustavo/projeto-musicoteca/internal/shared"
	"github.com/kustavo/projeto-musicoteca/internal/sync"
)

type params struct {
	sourcePath      string
	destinationPath string
	extSource       string
	extDestination  string
	tagsFilter      []string
	artist          string
	artistFolder    bool
	printMap        bool
}

func main() {
	start := time.Now()
	log.Println("Executando: ", start.Format("2006-01-02 15:04:05"))

	params, err := parseFlags()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = shared.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Mapeando...")
	artists, err := mapping.Mapping(params.sourcePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Mapeado!")

	if params.destinationPath == "" {
		log.Println("Destino não informado, transferência não será realizada")
		return
	}

	medias := []shared.MediaDTO{}

	if shared.IndexOf(shared.Conf.Ext.Video, params.extSource) != -1 {
		filteredMusicVideos := mapping.FilterMusicVideosByArtist(artists, params.artist)
		filteredMusicVideos = mapping.FilterMusicVideosByExtension(filteredMusicVideos, params.extSource)
		filteredMusicVideos = mapping.FilterMusicVideosByTags(filteredMusicVideos, params.tagsFilter)
		medias = shared.MapVideoMusicsToMediaDTO(filteredMusicVideos)
	}

	if shared.IndexOf(shared.Conf.Ext.Audio, params.extSource) != -1 {
		filteredSongs := mapping.FilterSongsByArtist(artists, params.artist)
		filteredSongs = mapping.FilterSongsByExtension(filteredSongs, params.extSource)
		filteredSongs = mapping.FilterSongsByTags(filteredSongs, params.tagsFilter)
		medias = shared.MapSongsToMediaDTO(filteredSongs)
	}

	log.Println("Transferindo...")
	err = sync.Sync(medias, params.destinationPath, params.extDestination, params.artistFolder)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Transferido!")

	if params.printMap {
		mapping.PrintMap(artists)
	}

	log.Println("Tempo de execução: ", time.Since(start))
}

func parseFlags() (params, error) {
	p := params{}

	flag.StringVar(&p.sourcePath, "source", "", "Diretório raiz dos arquivos de origem")
	flag.StringVar(&p.destinationPath, "destination", "", "Diretório raiz dos arquivos de destino")
	extSource := flag.String("ext-source", "flac", "Filtrar por extensão de arquivo na origem")
	extDestination := flag.String("ext-destination", "mp3", "Formato para o qual o arquivo será convertido")
	tagsFilter := flag.String("tags", "TOP", "Filtrar por arquivos que possuem todas as tags, separadas por vírgula")
	flag.StringVar(&p.artist, "artist", "", "Transferir arquivos de um artista específico")
	artistFolder := flag.String("artist-folder", "false", "Criar uma pasta separada para cada artista")
	printMap := flag.String("print-map", "false", "Imprimir o mapeamento dos arquivos")

	flag.Parse()

	if p.sourcePath == "" {
		return p, fmt.Errorf("local da origem não informado")
	}

	p.artistFolder = *artistFolder == "true"
	p.printMap = *printMap == "true"
	p.extSource = shared.NormalizeExtension(*extSource)
	p.extDestination = shared.NormalizeExtension(*extDestination)
	p.tagsFilter = strings.Split(*tagsFilter, ",")

	return p, nil
}
