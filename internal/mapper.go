package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	ArtistsFolder = "artistsFolder"
	ArtistFolder  = "artistFolder"
	AlbumFolder   = "albumFolder"
)

const (
	AudioMediaType = "audio"
	VideoMediaType = "video"
)

var (
	AudioMediaExtensions = []string{"flac", "mp3"}
	VideoMediaExtensions = []string{"mp4", "webm"}
)

var IgnoredArtistsFolderFiles = []string{"Readme.md", ".sync.ffs_db"}

type Song struct {
	artist    string
	name      string
	disc      int
	track     int
	flags     map[string]bool
	path      string
	midiaType string
}

type Album struct {
	artist string
	name   string
	flags  map[string]bool
	songs  map[string]Song
	path   string
}

type Artist struct {
	name   string
	albums map[string]Album
	path   string
}

func GetFilesMap(rootDir string, sourceDir string) (map[string]Artist, error) {
	songs := make(map[string]Artist)

	err := filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			isEmpty, err := isDirectoryEmpty(path)
			if err != nil {
				return err
			}
			if isEmpty {
				return fmt.Errorf("empty directory found: %s", path)
			}
		} else {
			err := appendSong(songs, path, rootDir)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = validateArtistName(songs)
	if err != nil {
		return nil, err
	}

	err = validateDiscoTrack(songs)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func appendFileMap(m map[string]Artist, newArtist Artist, newAlbum Album, newSong Song) {
	artistObj, artistExists := m[newArtist.name]
	if !artistExists {
		newArtist.albums = make(map[string]Album)
		artistObj = newArtist
		m[newArtist.name] = artistObj
	}

	albumMap := artistObj.albums
	albumObj, albumExists := albumMap[newAlbum.name]
	if !albumExists {
		newAlbum.songs = make(map[string]Song)
		albumObj = newAlbum
		albumMap[newAlbum.name] = albumObj
	}

	songMap := albumObj.songs
	// songMapKey := newSong.path
	// Irá acusar se tiver nomes iguais em tracks diferentes
	songMapKey := fmt.Sprintf("%s (%s)", newSong.name, newSong.midiaType)
	if _, songExists := songMap[songMapKey]; !songExists {
		songMap[songMapKey] = newSong
	}
}

func appendSong(songs map[string]Artist, filePath string, rootDir string) error {
	songObj := Song{}
	albumObj := Album{}
	artistObj := Artist{}

	relFilePath, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		return err
	}

	relativePathArr := strings.Split(relFilePath, string(filepath.Separator))

	if relativePathArr[0] == "_000_analise" || relativePathArr[0] == "_000_fila" { // TODO: remover
		return nil
	}

	if strings.Contains(relFilePath, "  ") {
		return fmt.Errorf("two spaces detected: %s", filePath)
	}

	isDirArtists := len(relativePathArr) == 1
	isDirArtist := len(relativePathArr) == 2
	isDirAlbum := len(relativePathArr) == 3

	artistPath := ""
	albumPath := ""
	songPath := ""

	if isDirArtists {
		if !contains(IgnoredArtistsFolderFiles, relativePathArr[0]) {
			return fmt.Errorf("invalid file location: %s", filePath)
		}
		return nil
	} else if isDirArtist {
		artistPath = filepath.Dir(filePath)
		albumPath = ""
		songPath = relFilePath
	} else if isDirAlbum {
		artistPath = filepath.Dir(filepath.Dir(filePath))
		albumPath = filepath.Dir(relFilePath)
		songPath = relFilePath
	} else {
		return fmt.Errorf("directory beyond depth limit: %s", filePath)
	}

	artistObj.path = artistPath
	artistObj.name = relativePathArr[0]
	if albumPath == "" {
		albumObj.path = ""
		albumObj.name = ""
	} else {
		albumObj.path = albumPath
		albumObj.artist, err = getArtist(relativePathArr[1])
		if err != nil {
			return fmt.Errorf("%s: %s", err.Error(), albumObj.path)
		}
		albumObj.name, err = getTitle(relativePathArr[1], false)
		if err != nil {
			return err
		}
		albumObj.flags, err = getFlags(relativePathArr[1])
		if err != nil {
			return err
		}
		songObj.disc, songObj.track, err = getDiscTrack(relativePathArr[2])
		if err != nil {
			return err
		}
	}
	songObj.path = songPath
	songObj.artist, err = getArtist(relativePathArr[len(relativePathArr)-1])
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), songObj.path)
	}
	songObj.name, err = getTitle(relativePathArr[len(relativePathArr)-1], true)
	if err != nil {
		return err
	}
	songObj.flags, err = getFlags(relativePathArr[len(relativePathArr)-1])
	if err != nil {
		return err
	}
	songObj.midiaType, err = getMidiaType(relativePathArr[len(relativePathArr)-1])
	if err != nil {
		return err
	}

	if len(albumObj.flags) > 0 && len(songObj.flags) > 0 {
		return fmt.Errorf("album and song have flags: %s", filePath)
	}

	if len(albumObj.flags) == 0 && len(songObj.flags) == 0 {
		return fmt.Errorf("album and song have no flags: %s", filePath)
	}

	appendFileMap(songs, artistObj, albumObj, songObj)

	return nil
}

func getArtist(fileName string) (string, error) {
	firstIndex := 0
	lastIndex := strings.Index(fileName, " - ")
	if lastIndex < 0 {
		return "", fmt.Errorf("artist and title separetor not found")
	}
	return fileName[firstIndex:lastIndex], nil
}

func getTitle(fileName string, isFile bool) (string, error) {
	firstIndex := strings.LastIndex(fileName, " - ") + 3
	lastIndex := strings.Index(fileName, "[")

	if lastIndex > 0 {
		if fileName[lastIndex-1:lastIndex] != " " {
			return "", fmt.Errorf("need a space between name and flags: %s", fileName)
		}
		lastIndex = lastIndex - 1
	} else {
		if isFile {
			lastIndex = strings.LastIndex(fileName, ".")
		} else {
			lastIndex = len(fileName)
		}
	}
	return fileName[firstIndex:lastIndex], nil
}

func getDiscTrack(fileName string) (int, int, error) {
	splFileName := strings.Split(fileName, " - ")
	rawTrack := ""
	if len(splFileName) > 2 {
		rawTrack = splFileName[1]
	}

	reg := regexp.MustCompile(`^([0-9]{2}).([0-9]{2})$`)
	matches := reg.FindStringSubmatch(rawTrack)

	if len(matches) == 3 {
		disc, _ := strconv.Atoi(matches[1])
		track, _ := strconv.Atoi(matches[2])
		return disc, track, nil
	}
	return 0, 0, nil
}

func getFlags(fileName string) (map[string]bool, error) {
	flags := map[string]bool{}

	// Obter substring das flags
	firstIndex := strings.Index(fileName, "[")
	lastIndex := strings.LastIndex(fileName, "]")

	if firstIndex < 0 {
		return flags, nil
	}

	flagsRaw := fileName[firstIndex : lastIndex+1]

	// Verificar se possui pares de colchetes e se não tem espaços
	if strings.Count(flagsRaw, "[") != strings.Count(flagsRaw, "]") || strings.Count(flagsRaw, " ") > 0 {
		return flags, fmt.Errorf("incorrectly formatted flags: %s", fileName)
	}

	// Converter para ; para depois realizar split
	flagsRaw = strings.Replace(flagsRaw, "][", ";", -1)
	flagsRaw = strings.Replace(flagsRaw, "]", "", -1)
	flagsRaw = strings.Replace(flagsRaw, "[", "", -1)
	flagsArr := strings.Split(flagsRaw, ";")

	for _, flag := range flagsArr {
		flags[flag] = true
	}

	return flags, nil
}

func getMidiaType(fileName string) (string, error) {
	firstIndex := strings.LastIndex(fileName, ".")
	if firstIndex < 0 {
		return "", fmt.Errorf("no extension found: %s", fileName)
	}
	extension := fileName[firstIndex+1:]

	if contains(AudioMediaExtensions, extension) {
		return AudioMediaType, nil
	} else if contains(VideoMediaExtensions, extension) {
		return VideoMediaType, nil
	} else {
		return "", fmt.Errorf("extension not allowed: %s", fileName)
	}
}

func validateArtistName(songs map[string]Artist) error {
	splitChar := " & "
	for _, artist := range songs {
		artistsName := strings.Split(artist.name, splitChar)
		for _, album := range artist.albums {
			if album.name != "" {
				artistsAlbum := strings.Split(album.artist, splitChar)
				for _, artistName := range artistsName {
					if !contains(artistsAlbum, artistName) {
						return fmt.Errorf("different album artist name: '%s' <> '%s'", artist.name, album.artist)
					}
				}
			}
			for _, song := range album.songs {
				artistsSong := strings.Split(song.artist, splitChar)
				for _, artistName := range artistsName {
					if !contains(artistsSong, artistName) {
						return fmt.Errorf("different song artist name: '%s' <> '%s'", artist.name, song.artist)
					}
				}
			}
		}
	}
	return nil
}

func validateDiscoTrack(songs map[string]Artist) error {
	for _, artist := range songs {
		for _, album := range artist.albums {
			if album.name != "" {
				discosTracks := map[int]map[int]bool{}
				for _, song := range album.songs {
					if song.midiaType != VideoMediaType {
						if _, ok := discosTracks[song.disc]; !ok {
							discosTracks[song.disc] = map[int]bool{}
						}
						discosTracks[song.disc][song.track] = true
					}
				}
				for j := 1; j <= len(discosTracks); j++ {
					if _, ok := discosTracks[j]; !ok {
						return fmt.Errorf("disc '%d' of album '%s' not found", j, album.path)
					}
					for i := 1; i <= len(discosTracks[j]); i++ {
						if _, ok := discosTracks[j][i]; !ok {
							return fmt.Errorf("track '%d.%d' of album '%s' not found", j, i, album.path)
						}
					}
				}
			}
		}
	}

	return nil
}
