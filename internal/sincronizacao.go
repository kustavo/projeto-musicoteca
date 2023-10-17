package internal

import (
	"path/filepath"
	"strings"
)

func ObterArquivosRemocao(mapOrigem map[string]Artista, mapDestino map[string]Artista) []string {
	itens := []string{}

	for keyBanda, bandaOrigem := range mapOrigem {
		if bandaDestino, ok := mapDestino[keyBanda]; !ok {
			itens = append(itens, bandaOrigem.path)
		} else {
			for keyAlbum, albumOrigem := range bandaOrigem.albuns {
				if albumDestino, ok := bandaDestino.albuns[keyAlbum]; !ok {
					itens = append(itens, albumOrigem.path)
				} else {
					for keyMusica, musicaOrigem := range albumOrigem.musicas {
						extensao := filepath.Ext(musicaOrigem.path)
						keyMusicaConvertida := strings.TrimSuffix(keyMusica, extensao) + obterExtensaoOriginal(extensao)

						if _, ok := albumDestino.musicas[keyMusicaConvertida]; !ok {
							itens = append(itens, musicaOrigem.path)
						}
					}
				}
			}
		}
	}
	return itens
}

func ObterArquivosInclusao(mapOrigem map[string]Artista, mapDestino map[string]Artista) []string {
	itens := []string{}

	for keyBanda, bandaOrigem := range mapOrigem {
		for keyAlbum, albumOrigem := range bandaOrigem.albuns {
			for keyMusica, musicaOrigem := range albumOrigem.musicas {
				extensao := filepath.Ext(musicaOrigem.path)
				keyMusicaConvertida := strings.TrimSuffix(keyMusica, extensao) + obterExtensaoConvertida(extensao)

				if bandaDestino, ok := mapDestino[keyBanda]; !ok {
					itens = append(itens, musicaOrigem.path)
				} else if albumDestino, ok := bandaDestino.albuns[keyAlbum]; !ok {
					itens = append(itens, musicaOrigem.path)
				} else if _, ok := albumDestino.musicas[keyMusicaConvertida]; !ok {
					itens = append(itens, musicaOrigem.path)
				}
			}
		}
	}

	return itens
}
