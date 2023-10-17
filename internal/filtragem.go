package internal

import (
	"path/filepath"
	"strings"
)

func Filtrar(mapa map[string]Artista, somenteAudio bool, somenteVideo bool, flagsStr string) error {
	flags := []string{}
	if flagsStr != "" {
		flags = strings.Split(flagsStr, ",")
	}

	for _, artista := range mapa {
		for _, album := range artista.albuns {
			for _, musica := range album.musicas {
				keyMusica := filepath.Base(musica.path)
				if somenteAudio && musica.tipo != TipoMidiaAudio || somenteVideo && musica.tipo != TipoMidiaVideo {
					delete(album.musicas, keyMusica)
				}
				var flagsMusica map[string]bool
				if len(musica.flags) > 0 {
					flagsMusica = musica.flags
				} else {
					flagsMusica = album.flags
				}
				for _, flag := range flags {
					if _, ok := flagsMusica[flag]; !ok {
						delete(album.musicas, keyMusica)
					}
				}
			}
		}
	}

	return nil
}
