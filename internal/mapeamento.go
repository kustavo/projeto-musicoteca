package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Musica struct {
	nome  string
	disco int
	faixa int
	flags map[string]bool
	path  string
	tipo  string
	album *Album
}

type Album struct {
	nome    string
	flags   map[string]bool
	musicas map[string]Musica
	path    string
	artista *Artista
}

type Artista struct {
	nome   string
	albuns map[string]Album
	path   string
}

func Mapear(path string, nomeArtista string) (map[string]Artista, error) {
	artistas, err := obterArtistas(path)
	if err != nil {
		return artistas, nil
	}

	for keyArtista := range artistas {
		artista := artistas[keyArtista]
		if nomeArtista == "" || artista.nome == nomeArtista {
			if err != nil {
				return artistas, nil
			}
			albuns, err := obterAlbuns(artista)
			if err != nil {
				return artistas, fmt.Errorf("[Artista: %s] %s", artista.nome, err)
			}
			artista.albuns = albuns
			artistas[keyArtista] = artista

			for keyAlbum := range artista.albuns {
				album := artista.albuns[keyAlbum]
				musicas, err := obterMusicas(album)
				if err != nil {
					return artistas, fmt.Errorf("[Artista: %s] [Album: %s] %s", artista.nome, album.nome, err)
				}
				album.musicas = musicas
				artista.albuns[keyAlbum] = album
			}
		}
	}
	return artistas, nil
}

func obterArtistas(path string) (map[string]Artista, error) {
	artistas := make(map[string]Artista)

	nomesDir, err := obterDiretorios(path)
	if err != nil {
		return artistas, nil
	}

	for _, nomeDir := range nomesDir {
		artistas[nomeDir] = Artista{
			nome: nomeDir,
			path: filepath.Join(path, nomeDir),
		}
	}
	return artistas, nil
}

func obterAlbuns(artista Artista) (map[string]Album, error) {
	albuns := make(map[string]Album)
	albuns[""] = Album{path: artista.path, artista: &artista}

	nomesDir, err := obterDiretorios(artista.path)
	if err != nil {
		return albuns, nil
	}

	for _, nomeDir := range nomesDir {
		artistaAlbum, err := obterArtistaArquivo(nomeDir)
		if err != nil {
			return albuns, fmt.Errorf("[Album: %s] %s", nomeDir, err)
		}
		if artistaAlbum != artista.nome {
			artistaAlbumArr := strings.Split(artistaAlbum, " & ")
			if !contains(artistaAlbumArr, artista.nome) {
				return albuns, fmt.Errorf("[Album: %s] %s", nomeDir, "nome do artista no album não corresponde")
			}
		}

		nomeAlbum, err := obterTituloArquivoAlbum(nomeDir)
		if err != nil {
			return albuns, fmt.Errorf("[Album: %s] %s", nomeDir, err)
		}

		flagsAlbum, err := obterFlagsArquivo(nomeDir)
		if err != nil {
			return albuns, fmt.Errorf("[Album: %s] %s", nomeDir, err)
		}

		albuns[nomeDir] = Album{
			nome:    nomeAlbum,
			flags:   flagsAlbum,
			path:    filepath.Join(artista.path, nomeDir),
			artista: &artista,
		}
	}
	return albuns, nil
}

func obterMusicas(album Album) (map[string]Musica, error) {
	musicas := make(map[string]Musica)

	isEmpty, err := isDiretorioVazio(album.path)
	if err != nil {
		return musicas, err
	}
	if isEmpty {
		return musicas, fmt.Errorf("diretório vazio: %s", album.path)
	}

	entries, err := os.ReadDir(album.path)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		nomeDir := e.Name()

		if e.IsDir() {
			if album.nome != "" {
				return musicas, fmt.Errorf("sub diretório não permitido: %s", filepath.Join(album.path, nomeDir))
			}
		} else {
			artistaMusica, err := obterArtistaArquivo(nomeDir)
			if err != nil {
				return musicas, fmt.Errorf("[Musica: %s] %s", nomeDir, err)
			}
			if artistaMusica != album.artista.nome {
				artistaMusicaArr := strings.Split(artistaMusica, " & ")
				if !contains(artistaMusicaArr, album.artista.nome) {
					return musicas, fmt.Errorf("[Musica: %s] nome do artista na música não corresponde", nomeDir)
				}
			}

			nomeMusica, err := obterTituloArquivoMusica(nomeDir)
			if err != nil {
				return musicas, fmt.Errorf("[Musica: %s] %s", nomeDir, err)
			}

			flagsMusica, err := obterFlagsArquivo(nomeDir)
			if err != nil {
				return musicas, fmt.Errorf("[Musica: %s] %s", nomeDir, err)
			}
			if album.nome != "" {
				if len(album.flags) > 0 && len(flagsMusica) > 0 {
					return musicas, fmt.Errorf("[Musica: %s] album e música não possuem flags", nomeDir)
				}
				if len(album.flags) == 0 && len(flagsMusica) == 0 {
					return musicas, fmt.Errorf("[Musica: %s] album e música não possuem flags", nomeDir)
				}
			} else {
				if len(flagsMusica) == 0 {
					return musicas, fmt.Errorf("[Musica: %s] música não possue flags", nomeDir)
				}
			}

			var discoMusica, faixaMusica int
			if album.nome != "" {
				discoMusica, faixaMusica, err = obterDiscoFaixa(nomeDir)
				if err != nil {
					return musicas, fmt.Errorf("[Musica: %s] %s", nomeDir, err)
				}
			}

			tipoMidia, err := obterTipoMidia(nomeDir)
			if err != nil {
				return musicas, fmt.Errorf("[Musica: %s] %s", nomeDir, err)
			}

			musica := Musica{
				nome:  nomeMusica,
				flags: flagsMusica,
				path:  filepath.Join(album.path, nomeDir),
				disco: discoMusica,
				faixa: faixaMusica,
				tipo:  tipoMidia,
				album: &album,
			}
			musicas[nomeDir] = musica
		}
	}

	if album.nome != "" {
		err = validarDiscoFaixa(musicas)
		if err != nil {
			return musicas, fmt.Errorf("[Musica: %s] %s", album.nome, err)
		}
	}

	return musicas, nil
}

func obterTituloArquivoAlbum(nomeDir string) (string, error) {
	firstIndex := strings.LastIndex(nomeDir, " - ") + 3
	lastIndex := strings.Index(nomeDir, "[")

	if lastIndex > 0 {
		if nomeDir[lastIndex-1:lastIndex] != " " {
			return "", fmt.Errorf("espaço entre nome do album e flag ausente: %s", nomeDir)
		}
		lastIndex = lastIndex - 1
	} else {
		lastIndex = len(nomeDir)
	}
	return nomeDir[firstIndex:lastIndex], nil
}

func obterArtistaArquivo(nomeDir string) (string, error) {
	firstIndex := 0
	lastIndex := strings.Index(nomeDir, " - ")
	if lastIndex < 0 {
		return "", fmt.Errorf("separador entre o nome da banda e titulo incorreto: '%s'", nomeDir)
	}
	return nomeDir[firstIndex:lastIndex], nil
}

func obterFlagsArquivo(nomeArquivo string) (map[string]bool, error) {
	flags := map[string]bool{}

	firstIndex := strings.Index(nomeArquivo, "[")
	lastIndex := strings.LastIndex(nomeArquivo, "]")

	if firstIndex < 0 {
		return flags, nil
	}

	flagsRaw := nomeArquivo[firstIndex : lastIndex+1]

	// Verificar se possui pares de colchetes e se não tem espaços
	if strings.Count(flagsRaw, "[") != strings.Count(flagsRaw, "]") || strings.Count(flagsRaw, " ") > 0 {
		return flags, fmt.Errorf("flag formatada de forma incorreta: %s", nomeArquivo)
	}

	// Converter para ; para depois realizar split
	flagsRaw = strings.Replace(flagsRaw, "][", ";", -1)
	flagsRaw = strings.Replace(flagsRaw, "]", "", -1)
	flagsRaw = strings.Replace(flagsRaw, "[", "", -1)
	flagsArr := strings.Split(flagsRaw, ";")

	for i := 0; i < len(flagsArr); i++ {
		flag := flagsArr[i]
		if i == 0 {
			if !contains(MainFlags, flag) {
				return flags, fmt.Errorf("flag princiapl '%s' não encontrada: %s", flag, nomeArquivo)
			}
		} else {
			if !contains(SecondaryFlags, flag) {
				return flags, fmt.Errorf("flag secundária '%s' não encontrada: %s", flag, nomeArquivo)
			}
		}
		_, existe := flags[flag]
		if existe {
			return flags, fmt.Errorf("flag duplicada '%s': %s", flag, nomeArquivo)
		}
		flags[flag] = true
	}

	return flags, nil
}

func obterTituloArquivoMusica(nomeDir string) (string, error) {
	nomeArquivo := filepath.Base(nomeDir)

	firstIndex := strings.LastIndex(nomeArquivo, " - ") + 3
	lastIndex := strings.Index(nomeArquivo, "[")

	if lastIndex > 0 {
		if nomeArquivo[lastIndex-1:lastIndex] != " " {
			return "", fmt.Errorf("espaço entre nome da música e flag ausente: %s", nomeDir)
		}
		lastIndex = lastIndex - 1
	} else {
		lastIndex = strings.LastIndex(nomeArquivo, ".")
	}
	return nomeArquivo[firstIndex:lastIndex], nil
}

func obterDiscoFaixa(nomeArquivo string) (int, int, error) {
	splnomeArquivo := strings.Split(nomeArquivo, " - ")
	rawTrack := ""
	if len(splnomeArquivo) > 2 {
		rawTrack = splnomeArquivo[1]
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

func obterTipoMidia(nomeArquivo string) (string, error) {
	firstIndex := strings.LastIndex(nomeArquivo, ".")
	if firstIndex < 0 {
		return "", fmt.Errorf("extensão não encontrada: %s", nomeArquivo)
	}
	extensao := nomeArquivo[firstIndex+1:]

	if contains(ExtensoesMidiasAudio, extensao) {
		return TipoMidiaAudio, nil
	} else if contains(ExtensoesMidiasVideo, extensao) {
		return TipoMidiaVideo, nil
	} else {
		return "", fmt.Errorf("extensão não permitida: %s", extensao)
	}
}

func validarDiscoFaixa(musicas map[string]Musica) error {
	faixasDiscos := map[int]map[int]bool{}
	for _, musica := range musicas {
		if musica.tipo != TipoMidiaVideo {
			if _, ok := faixasDiscos[musica.disco]; !ok {
				faixasDiscos[musica.disco] = map[int]bool{}
			}
			faixasDiscos[musica.disco][musica.faixa] = true
		}
	}
	for j := 1; j <= len(faixasDiscos); j++ {
		if _, ok := faixasDiscos[j]; !ok {
			return fmt.Errorf("disco '%d' não encontrado", j)
		}
		for i := 1; i <= len(faixasDiscos[j]); i++ {
			if _, ok := faixasDiscos[j][i]; !ok {
				return fmt.Errorf("faixa '%d.%d' não encontrada", j, i)
			}
		}
	}

	return nil
}
