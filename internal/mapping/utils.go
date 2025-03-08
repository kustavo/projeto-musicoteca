package mapping

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kustavo/projeto-musicoteca/internal/shared"
)

func GetTags(str string) []string {
	tags := []string{}

	start := strings.Index(str, " [")
	if start == -1 {
		return tags
	}

	substr := str[start+2:]
	substr = substr[:len(substr)-1]
	tags = strings.Split(substr, "][")

	return tags
}

func GetDiscAndTrack(str string) (int, int) {
	parts := strings.Split(str, " - ")
	discAndTrack := strings.Split(parts[1], ".")

	disc, err := strconv.Atoi(discAndTrack[0])
	if err != nil {
		return 0, 0
	}
	track, err := strconv.Atoi(discAndTrack[1])
	if err != nil {
		return 0, 0
	}
	return disc, track
}

func GetDirCaptions(path string) []string {
	captions := []string{}

	entries, err := os.ReadDir(path)
	if err != nil {
		return captions
	}

	for _, e := range entries {
		extension := filepath.Ext(e.Name())
		if !e.IsDir() && shared.IndexOf(shared.Conf.Ext.Caption, extension) != -1 {
			captions = append(captions, e.Name())
		}
	}

	return captions
}

func GetCaptions(name string, path string) []string {
	list := []string{}
	captions := GetDirCaptions(path)
	hasMainCaption := false

	for _, caption := range captions {
		extension := filepath.Ext(caption)
		capName := strings.TrimSuffix(caption, extension)
		if name == capName {
			list = append(list, caption)
			hasMainCaption = true
		} else {
			if strings.HasPrefix(capName, name) {
				tags := GetTags(capName)
				if shared.IndexOf(shared.Conf.Tags.Language, tags[len(tags)-1]) != -1 {
					list = append(list, caption)
				}
			}
		}

	}

	if !hasMainCaption && len(captions) > 0 {
		return []string{}
	}

	return list
}

func GetUnusedCaptions(path string, validCaptions []string) []string {
	list := []string{}
	dirCaptions := GetDirCaptions(path)
	for _, c := range dirCaptions {
		if shared.IndexOf(validCaptions, c) == -1 {
			list = append(list, c)
		}
	}
	return list
}

func FilterSongsByArtist(artists []shared.Artist, name string) []shared.Song {
	var filtered []shared.Song
	for _, artist := range artists {
		if name == "" || artist.Name == name {
			for _, album := range artist.Albums {
				filtered = append(filtered, album.Songs...)
			}
		}
	}
	return filtered
}

func FilterMusicVideosByArtist(artists []shared.Artist, name string) []shared.MusicVideo {
	var filtered []shared.MusicVideo
	for _, artist := range artists {
		if name == "" || artist.Name == name {
			filtered = append(filtered, artist.MusicVideos...)
		}
	}
	return filtered
}

func FilterSongsByExtension(songs []shared.Song, extension string) []shared.Song {
	var filtered []shared.Song
	for _, song := range songs {
		if song.Extension == extension {
			filtered = append(filtered, song)
		}
	}

	return filtered
}

func FilterMusicVideosByExtension(musicVideos []shared.MusicVideo, extension string) []shared.MusicVideo {
	var filtered []shared.MusicVideo
	for _, musicVideo := range musicVideos {
		if musicVideo.Extension == extension {
			filtered = append(filtered, musicVideo)
		}
	}

	return filtered
}

func FilterSongsByTags(songs []shared.Song, tags []string) []shared.Song {
	var filtered []shared.Song
	for _, song := range songs {
		if containsAllTags(song.Tags, tags) {
			filtered = append(filtered, song)
		}
	}
	return filtered
}

func FilterMusicVideosByTags(musicVideos []shared.MusicVideo, tags []string) []shared.MusicVideo {
	var filtered []shared.MusicVideo
	for _, musicVideo := range musicVideos {
		if containsAllTags(musicVideo.Tags, tags) {
			filtered = append(filtered, musicVideo)
		}
	}
	return filtered
}

func containsAllTags(itemTags, filterTags []string) bool {
	tagSet := make(map[string]struct{}, len(itemTags))
	for _, tag := range itemTags {
		tagSet[strings.ToUpper(tag)] = struct{}{}
	}

	for _, tag := range filterTags {
		if _, found := tagSet[strings.ToUpper(tag)]; !found {
			return false
		}
	}
	return true
}

func PrintMap(artists []shared.Artist) {
	for _, a := range artists {
		log.Println(a.Name)
		for _, al := range a.Albums {
			albumName := al.Name
			if al.Name == "" {
				albumName = "*"
			}
			log.Println("  ", albumName)
			for _, s := range al.Songs {
				log.Println("    ", s.Name)
			}
		}
	}
}
