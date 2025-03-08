package mapping

import (
	"errors"
	"strings"

	"github.com/kustavo/projeto-musicoteca/internal/shared"
)

func Mapping(path string) ([]shared.Artist, error) {
	artists, err := mappingArtists(path)
	if err != nil {
		return nil, err
	}

	for i := range artists {
		musicVideos, err := mappingMusicVideos(artists[i])
		if err != nil {
			return nil, err
		}
		artists[i].MusicVideos = musicVideos

		albums, err := mappingAlbums(artists[i])
		if err != nil {
			return nil, err
		}

		artists[i].Albums = albums

		for j := range albums {
			songs, err := mappingSongs(albums[j])
			if err != nil {
				return nil, err
			}
			albums[j].Songs = songs
		}
	}

	return artists, nil
}

func mappingArtists(path string) ([]shared.Artist, error) {
	list := []shared.Artist{}

	files, err := shared.ListDirectoryFiles(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {

		if !IsDir(f) {
			return nil, errors.New("O arquivo não é um diretório: " + f.Name)
		}

		if !IsSpacesFormatValid(f) {
			return nil, errors.New("O arquivo possui espaços extras: " + f.Name)
		}

		artist := shared.Artist{Name: f.Name, Path: f.Path}
		list = append(list, artist)
	}

	return list, nil
}

func mappingAlbums(artist shared.Artist) ([]shared.Album, error) {
	list := []shared.Album{}

	defaultAlbum := shared.Album{Name: "", Path: artist.Path, Artist: &artist}
	list = append(list, defaultAlbum)

	files, err := shared.ListDirectoryFiles(artist.Path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {

		if !IsDir(f) {
			continue
		}

		if !IsSpacesFormatValid(f) {
			return nil, errors.New("O arquivo possui espaços extras: " + f.Name)
		}

		if !IsCommaFormatValid(f.Name) {
			return nil, errors.New("Vírgula fora do padrão: " + f.Name)
		}

		if !IsMatchWithRegex(f.Name, "^"+artist.Name+" - ") &&
			!IsMatchWithRegex(f.Name, "^.{2,}, "+artist.Name+" - ") &&
			!IsMatchWithRegex(f.Name, "^"+artist.Name+", .{2,} - ") {
			return nil, errors.New("Nome artista inválido: " + f.Name)
		}

		if !IsTagFormatValid(f.Name) {
			return nil, errors.New("Formato de tags inválida: " + f.Name)
		}

		tags := GetTags(f.Name)

		if !AreTagsValid(tags) {
			return nil, errors.New("Tags não existem ou fora de ordem: " + f.Name)
		}

		album := shared.Album{Name: f.Name, Path: f.Path, Tags: tags, Artist: &artist}
		list = append(list, album)
	}

	artist.Albums = list
	return artist.Albums, nil
}

func mappingSongs(album shared.Album) ([]shared.Song, error) {
	list := []shared.Song{}

	files, err := shared.ListDirectoryFiles(album.Path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {

		if IsDir(f) {
			continue
		}

		if !IsExtensionValid(f.Extension) {
			return nil, errors.New("Extensão desconhecida: " + f.Path)
		}

		if shared.IndexOf(shared.Conf.Ext.Audio, f.Extension) == -1 {
			continue
		}

		if !IsSpacesFormatValid(f) {
			return nil, errors.New("O arquivo possui espaços extras: " + f.Name)
		}

		if !IsCommaFormatValid(f.Name) {
			return nil, errors.New("Vírgula fora do padrão: " + f.Name)
		}

		if !IsMatchWithRegex(f.Name, "^"+album.Artist.Name+" - ") &&
			!IsMatchWithRegex(f.Name, "^.{2,}, "+album.Artist.Name+" - ") &&
			!IsMatchWithRegex(f.Name, "^"+album.Artist.Name+", .{2,} - ") {
			return nil, errors.New("Nome artista inválido: " + f.Name)
		}

		tags := []string{}
		track := 0
		disc := 0
		if album.Name != "" {

			if !IsTrackFormatValid(f.Name) {
				return nil, errors.New("Formato de faixa inválido: " + f.Name)
			}

			disc, track = GetDiscAndTrack(f.Name)

		} else {

			if !IsTagFormatValid(f.Name) {
				return nil, errors.New("Formato de tags inválida: " + f.Name)
			}

			tags = GetTags(f.Name)

			if !AreTagsValid(tags) {
				return nil, errors.New("Tags não existem ou fora de ordem: " + f.Name)
			}
		}

		song := shared.Song{Name: f.Name, Path: f.Path, Tags: tags, DiscNumber: disc, Track: track, Extension: f.Extension, Album: &album}
		list = append(list, song)
	}

	album.Songs = list

	if album.Name != "" {
		if !AreTrackOrderValid(album.Songs) {
			return nil, errors.New("As faixas do álbum não estão em ordem: " + album.Name)
		}
	}

	return album.Songs, nil
}

func mappingMusicVideos(artist shared.Artist) ([]shared.MusicVideo, error) {
	list := []shared.MusicVideo{}
	validCaptions := []string{}

	files, err := shared.ListDirectoryFiles(artist.Path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {

		if IsDir(f) {
			continue
		}

		if !IsExtensionValid(f.Extension) {
			return nil, errors.New("Extensão desconhecida: " + f.Path)
		}

		if shared.IndexOf(shared.Conf.Ext.Video, f.Extension) == -1 {
			continue
		}

		if !IsSpacesFormatValid(f) {
			return nil, errors.New("O arquivo possui espaços extras: " + f.Name)
		}

		if !IsCommaFormatValid(f.Name) {
			return nil, errors.New("Vírgula fora do padrão: " + f.Name)
		}

		if !IsMatchWithRegex(f.Name, "^"+artist.Name+" - ") &&
			!IsMatchWithRegex(f.Name, "^.{2,}, "+artist.Name+" - ") &&
			!IsMatchWithRegex(f.Name, "^"+artist.Name+", .{2,} - ") {
			return nil, errors.New("Nome artista inválido: " + f.Name)
		}

		if !IsTagFormatValid(f.Name) {
			return nil, errors.New("Formato de tags inválida: " + f.Name)
		}

		tags := GetTags(f.Name)

		if !AreTagsValid(tags) {
			return nil, errors.New("Tags não existem ou fora de ordem: " + f.Name)
		}

		captions := GetCaptions(f.Name, artist.Path)
		validCaptions = append(validCaptions, captions...)

		musicVideo := shared.MusicVideo{Name: f.Name, Path: f.Path, Tags: tags, Captions: captions, Extension: f.Extension, Artist: &artist}
		list = append(list, musicVideo)
	}

	unusedCaptions := GetUnusedCaptions(artist.Path, validCaptions)
	if len(unusedCaptions) > 0 {
		return nil, errors.New("Possui legendas sem referência: " + strings.Join(unusedCaptions, ", "))
	}

	return list, nil
}
