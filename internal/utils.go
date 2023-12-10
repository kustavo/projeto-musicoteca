// ## Tag principal

// * SPD - Speed Metal
// * HVY - Heavy Metal
// * HRD - Hard Rock
// * ROC - Rock
// * POP - Pop Rock (sem guitarras, refrão repetitivo)
// * SFT - Soft (rock lento ou romântico)
// * DSC - Disco (estilo anos 70)
// * DCE - Dance (estilo anos 90)
// * RLX - Relax (muito tranquila com voz) (para relaxamento, igual foco mas pode ter voz)
// * FOC - focus (muito tranquila sem voz) (para manter o foco, não distrair com barulho externo)
// * SPN - Spanish (músicas em espanhol)
// * INS - Instrumental
// * COU - Country
// * REG - Reggae

// ## Tags secundárias

// * 60S - Anos 60
// * 70S - Anos 70
// * 80S - Anos 80
// * 90S - Anos 90
// * LIV - Live (ao vivo, com platéia)
// * ORC - Orchestrated (músicas orquestradas)
// * COR - Choir (coro)
// * BRA - Nacional
// * GOV - Good Vibes (alto astral)
// * TOP - Top (favoritas)
// * NUL - Null (apenas para o acervo, não deve entrar em playlist)

package internal

import (
	"os"
	"path/filepath"
)

const (
	TipoMidiaAudio = "audio"
	TipoMidiaVideo = "video"
	TipoLetra      = "letra"
)

var (
	ExtensoesMidiasAudio = []string{"flac", "mp3"}
	ExtensoesMidiasVideo = []string{"mp4", "webm"}
	ExtensoesLetra       = []string{"srt"}
)

var (
	MainFlags      = []string{"SPD", "HVY", "HRD", "ROC", "POP", "SFT", "DSC", "DCE", "RLX", "FOC", "SPN", "INS", "COU", "REG"}
	SecondaryFlags = []string{"60S", "70S", "80S", "90S", "LIV", "ORC", "COR", "BRA", "GOV", "TOP", "NUL"}
)

func contains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}

func isDiretorioVazio(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			isEmpty, err := isDiretorioVazio(filepath.Join(path, entry.Name()))
			if err != nil {
				return false, err
			}
			if !isEmpty {
				return false, nil
			}
		} else {
			return false, nil
		}
	}

	return true, nil
}

func obterDiretorios(path string) ([]string, error) {
	diretorios := []string{}

	entries, err := os.ReadDir(path)
	if err != nil {
		return diretorios, err
	}

	for _, e := range entries {
		if e.IsDir() {
			nomeDir := e.Name()
			if (len(nomeDir) < 5 || nomeDir[:5] != "_000_") && nomeDir != ".Trash-1000" {
				diretorios = append(diretorios, nomeDir)
			}
		}
	}
	return diretorios, nil
}

func obterExtensaoConvertida(extensao string) string {
	switch extensao {
	case ".flac":
		return ".mp3"
	default:
		return extensao
	}
}

func obterExtensaoOriginal(extensao string) string {
	switch extensao {
	case ".mp3":
		return ".flac"
	default:
		return extensao
	}
}

func removerExtensaoPath(path string) string {
	extensao := filepath.Ext(path)
	pathSemExtensao := path[0 : len(path)-len(extensao)]
	return pathSemExtensao
}
