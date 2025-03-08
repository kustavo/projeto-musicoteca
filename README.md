# projeto-musicoteca

Script para converter e sincronizar arquivos de audio e vídeo entre uma origem e destino.

## Como usar

Arquivo de configuração (deve estar no mesmo diretório):

- **config.json**

Parâmetros:

- **-source**: Path de origem. Raiz onde se encontra os diretorios dos artistas na origem.
- **-destination**: Path de destino. Raiz onde se encontra os diretorios dos artistas no destino.
- **-ext-source** `(padrão "flac")`: Extensão para filtrar arquivos na origem.
- **-ext-destination** `(padrão "mp3")`: Extensão para o qual o arquivo será convertido.
- **-tags** `(padrão "top")`: Filtrar por arquivos que possuem todas as tags, separadas por vírgula.
- **-artist** `(padrão "")`: Transferir arquivos de um artista específico.
- **-artist-folder** `(padrão "false")`: Criar uma pasta separada para cada artista.
- **-print-map** `(padrão "false")`: Imprimir o mapeamento dos arquivos.


**Exemplo**: 

./projeto-musicoteca -source "/home/gustavo/Músicas/" -destination "/home/gustavo/Músicas/mp3/" -ext-source "flac" -ext-destination "mp3" -artist "a-ha" -tags "ROC,TOP" -artist-folder "true"


## Tags principais

* SPD - Speed Metal
* HVY - Heavy Metal
* HRD - Hard Rock
* ROC - Rock
* POP - Pop Rock (sem guitarras, refrão repetitivo)
* SFT - Soft (rock lento ou romântico)
* DSC - Disco (estilo anos 70)
* DCE - Dance (estilo anos 90)
* RLX - Relax (muito tranquila com voz) (para relaxamento, igual foco mas pode ter voz)
* FOC - focus (muito tranquila sem voz) (para manter o foco, não distrair com barulho externo)
* SPN - Spanish (músicas em espanhol)
* INS - Instrumental
* COU - Country
* REG - Reggae

## Tags secundárias

* 60S - Anos 60
* 70S - Anos 70
* 80S - Anos 80
* 90S - Anos 90
* LIV - Live (ao vivo, com platéia)
* ORC - Orchestrated (músicas orquestradas)
* COR - Choir (coro)
* BRA - Nacional
* GOV - Good Vibes (alto astral)
* TOP - Top (favoritas)
* NUL - Null (apenas para o acervo, não deve entrar em playlist)






