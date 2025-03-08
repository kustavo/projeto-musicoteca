package mapping

import (
	"regexp"
	"strings"

	"github.com/kustavo/projeto-musicoteca/internal/shared"
)

func IsDir(file shared.File) bool {
	return file.IsDir
}

func IsSpacesFormatValid(file shared.File) bool {
	return !strings.Contains(file.Name, "  ") && strings.TrimSpace(file.Name) == file.Name
}

func IsMatchWithRegex(str string, reg string) bool {
	regex := regexp.MustCompile(reg)
	return regex.MatchString(str)
}

func IsCommaFormatValid(str string) bool {
	if IsMatchWithRegex(str, ",\\S") {
		return false
	}

	if IsMatchWithRegex(str, " ,") {
		return false
	}
	return true
}

func IsTagFormatValid(str string) bool {
	start := strings.Index(str, " [")
	if start == -1 {
		return false
	}
	substr := str[start+1:]
	re := regexp.MustCompile(`^(\[[^\[\]]*\])+$`)
	return re.MatchString(substr)
}

func AreTagsValid(tags []string) bool {
	if !checkMainFlag(tags[0]) {
		return false
	}
	return checkMultFlags(tags[1:])
}

func checkMainFlag(tag string) bool {
	return shared.IndexOf(shared.Conf.Tags.Main, tag) != -1
}

func checkMultFlags(tags []string) bool {
	lastIndex := -1
	for _, tag := range tags {
		if i := shared.IndexOf(shared.Conf.Tags.Mult, tag); i == -1 || i < lastIndex {
			return false
		} else {
			lastIndex = i
		}
	}
	return true
}

func IsTrackFormatValid(str string) bool {
	start := strings.Index(str, " - ")
	if start == -1 {
		return false
	}
	substr := str[start+3 : start+8]

	re := regexp.MustCompile(`\d{2}\.\d{2}`)
	return re.MatchString(substr)
}

func AreTrackOrderValid(songs []shared.Song) bool {
	faixasDiscos := map[int]map[int]bool{}

	for _, s := range songs {
		if _, ok := faixasDiscos[s.DiscNumber]; !ok {
			faixasDiscos[s.DiscNumber] = map[int]bool{}
		}
		faixasDiscos[s.DiscNumber][s.Track] = true
	}

	for j := 1; j <= len(faixasDiscos); j++ {
		if _, ok := faixasDiscos[j]; !ok {
			return false
		}
		for i := 1; i <= len(faixasDiscos[j]); i++ {
			if _, ok := faixasDiscos[j][i]; !ok {
				return false
			}
		}
	}

	return true
}

func IsExtensionValid(extension string) bool {
	return shared.IndexOf(shared.Conf.Ext.Audio, extension) != -1 ||
		shared.IndexOf(shared.Conf.Ext.Video, extension) != -1 ||
		shared.IndexOf(shared.Conf.Ext.Caption, extension) != -1
}
