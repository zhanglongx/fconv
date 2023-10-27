# fconv

fconv is a tool for quick converting of file encoding and format.

fconv can detect the encoding and format of a file and convert it manually or automatically. For the sake of simplicity, we need to assume two "mode" here, Windows mode, which means GB-10830 encoding, DOS format (CRLF), and Unix mode, which means UTF-8 encoding, Unix format (LF). I know this is not very accurate, but it does fulfill more than 80% of the generality needs.

In manual mode, you can specify whether the desired file converted mode is Windows mode or Unix mode. In automatic mode, the conversion between Windows mode and Unix mode is done automatically.



⚠️ fconv is designed for GB-10830 and UTF-8, and it is not recommended to use it for other encodings.

⚠️ file encoding conversion can be a very dangerous operation sometimes, and leading to the lose of data. So you are responsible for your own actions.

## Usage

```bash
    $ fconv [options] <file>...
```

- `-p` prints the encoding and format of the file. And will do nothing about changing the file.

- `-a`, `-w`, `-u` specifies the desired mode of the file. `-a` means automatic mode, `-w` means Windows mode, `-u` means Unix mode. 

    ⚠️ fconv always do in-place conversion, so please backup your files before using these options.

## Note

- Since GB-2312 encoding is a subset of GB-10830 encoding, GB-2312 encoding is also supported.

- When detecting the "mode" of a file, the file encoding always takes precedence over the file format. This means that if a file is encoded in GB-10830 and formatted for Unix (LF), it will be considered a 'Windows mode' file; similarly, a UTF-8, Windows (CRLF) file will be considered a 'Unix mode' file.

- utf-8 bom encoded files. fconv supports utf-8 bom encoded files when reading files, but when saving files, it will only save as utf-8 (w/o bom) encoded files.

- fconv only supports GB-10830 and utf-8 encoded file conversion. So, if a file is utf-8 encoded and contains characters other than ASCII characters and Chinese simplified characters, fconv will usually refuse to convert it. 

- Due to the GB-18030 encoding method, the encoding of Emoji characters is not supported. Therefore, if a file is UTF-8 encoded but contains Emoji characters, then fconv will also refuse to convert it.
