package shared

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	Name      string
	Path      string
	IsDir     bool
	Extension string
}

type Config struct {
	Tags struct {
		Main     []string `json:"main"`
		Mult     []string `json:"mult"`
		Language []string `json:"language"`
	} `json:"tags"`
	Ext struct {
		Audio   []string `json:"audio"`
		Video   []string `json:"video"`
		Caption []string `json:"caption"`
	} `json:"extensions"`
	Ignore []string `json:"ignore"`
}

type Artist struct {
	Name        string
	Albums      []Album
	MusicVideos []MusicVideo
	Path        string
}

type Album struct {
	Name   string
	Tags   []string
	Songs  []Song
	Path   string
	Artist *Artist
}

type Song struct {
	Name       string
	DiscNumber int
	Track      int
	Tags       []string
	Path       string
	Album      *Album
	Extension  string
}

type MusicVideo struct {
	Name      string
	Tags      []string
	Path      string
	Captions  []string
	Extension string
	Artist    *Artist
}

type MediaDTO struct {
	ArtistName string
	Name       string
	Path       string
	Extension  string
}

var (
	confPath = "config.json"
	Conf     Config
)

func IndexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

func LoadConfig() error {
	var config Config

	file, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	Conf = config
	return nil
}

func NormalizeExtension(extension string) string {
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return extension
}

func ListDirectoryFiles(path string) ([]File, error) {
	list := []File{}

	entries, err := os.ReadDir(path)
	if err != nil {
		return list, err
	}

	ignorePattern := regexp.MustCompile(strings.Join(Conf.Ignore, "|"))

	for _, e := range entries {
		if ignorePattern.MatchString(e.Name()) {
			continue
		}

		extension := filepath.Ext(e.Name())
		file := File{
			Name:      strings.TrimSuffix(e.Name(), extension),
			Path:      filepath.Join(path, e.Name()),
			IsDir:     e.IsDir(),
			Extension: extension,
		}
		list = append(list, file)
	}

	return list, nil
}

func FilterFilesByExtension(files []File, extension string) []File {
	var filtered []File
	for _, file := range files {
		if file.Extension == extension {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func MapVideoMusicsToMediaDTO(musicVideos []MusicVideo) []MediaDTO {
	var medias []MediaDTO
	for _, mv := range musicVideos {
		medias = append(medias, MediaDTO{
			ArtistName: mv.Artist.Name,
			Name:       mv.Name,
			Path:       mv.Path,
			Extension:  mv.Extension,
		})
	}
	return medias
}

func MapSongsToMediaDTO(songs []Song) []MediaDTO {
	var medias []MediaDTO
	for _, song := range songs {
		medias = append(medias, MediaDTO{
			ArtistName: song.Album.Artist.Name,
			Name:       song.Name,
			Path:       song.Path,
			Extension:  song.Extension,
		})
	}
	return medias
}
