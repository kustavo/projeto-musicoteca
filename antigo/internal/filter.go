package internal

import (
	"fmt"
)

func Filter(sourceFiles map[string]Artist, destinyFiles map[string]Artist, includeAudio bool, includeVideo bool, flags []string) error {
	fmt.Println("Filtering files...")

	for _, sourceArtist := range sourceFiles {
		for _, sourceAlbums := range sourceArtist.albums {
			for _, sourceSong := range sourceAlbums.songs {
				if !isTransferable(sourceSong, includeAudio, includeVideo, flags) {
					delete(sourceAlbums.songs, sourceSong.path)
				}

				destinyArtist := destinyFiles[sourceArtist.name]
				destinyAlbum := destinyArtist.albums[sourceAlbums.name]
				songMapKey := fmt.Sprintf("%s (%s)", sourceSong.name, sourceSong.midiaType)
				destinySong := destinyAlbum.songs[songMapKey]

				if destinySong.path != "" {
					delete(sourceAlbums.songs, songMapKey)
				}
			}
		}
	}

	return nil
}

func isTransferable(song Song, includeAudio bool, includeVideo bool, flags []string) bool {
	if includeAudio && song.midiaType == AudioMediaType || includeVideo && song.midiaType == VideoMediaType {
		for _, flag := range flags {
			if _, flagExists := song.flags[flag]; !flagExists {
				return false
			}
		}
	}
	return true
}
