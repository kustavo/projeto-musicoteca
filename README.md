# projeto-musicoteca

Script para converter e sincronizar arquivos de audio e vídeo entre uma origem e destino.

## Como usar

Parâmetros:

- **origem**: Path de origem. Raiz onde se encontra os diretorios dos artistas na origem.
- **destino**: Path de destino. Raiz onde se encontra os diretorios dos artistas no destino.
- **artista**: (Opcional) Especifica um artista para a transferência. Caso não seja informado será selecionado todos artistas.
- **somente-audio**: (Opcional) Será transferido somente arquivos de áudio.
- **somente-video**: (Opcional) Será transferido somente arquivos de vídeo.
- **flags**: (Opcional) Flags usadas na filtragem. Será usado o modo AND. Irá buscas músicas que tenha todas as flags passadas. As flags dever separadas por "," virgula. 

**Exemplo**: 

./projeto-musicoteca -origem "/mnt/arquivos/Músicas/" -destino "/mnt/usb/musicas/" -artista "Queen" -somente-audio -flags "ROC,TOP"  
