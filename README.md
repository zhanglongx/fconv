# fconv

fconv is a tool for quick converting of file encoding and format.

⚠️ fconv is designed for GB-2312 and UTF-8, and it is not recommended to use it for other encodings.

⚠️ file encoding conversion can be a very dangerous operation sometimes, and leading to the lose of data. So you are responsible for your own actions.

## Background

Converting between different file encodings is not very easy. It can be very difficult if we want to devise a comprehensive set of conversion rules. Therefore, it is necessary to simplify the problem. Currently, the Chinese character encoding standard is GB-18030. GB-18030 has many advantages, such as support for unicode base plane(BMP). but a reality problem here is that under common operating systems, such as Windows, which uses CP936 as default, there is no and full support for GB-18030 but GB-2312. Meanwhile, the most common 80%+ problem we have is that a file is created under Windows, but it doesn't open properly under Linux (or MacOS). So, in order to solve this 80% of scenarios, it is necessary for us to make fconv simpler as well.

fconv defines two "file modes": Windows mode and Unix mode, where files are considered to be GB-2312 encoded in DOS format (CRLF), and Unix mode, where files are considered to be UTF-8 encoded in Unix format (LF). fconv main goal is to convert between these two file modes. The reason for making this assumption is that this scenario is very common. And other scenarios are very uncommon. For example, a Windows user set CP54936 to create a file with GB-18030 encoding; or a Linux user uses GB-18030 as locale.

fconv has two "conversion modes": manual and automatic. In manual mode, the user specifies the desired "file mode" to be obtained, while in automatic mode, the conversion between the two file modes is performed automatically.

## Usage

```bash
    $ fconv [options] <file>...
```

- `-p` prints the encoding and format of the file. And will do nothing about changing the file.

- `-a`, `-w`, `-u` specifies the desired 'file mode'. `-a` means automatic mode, `-w` means Windows mode, `-u` means Unix mode. 

    ⚠️ fconv always do in-place conversion, so please backup your files before using these options.

## Note

- When detecting the "mode" of a file, the file encoding always takes precedence over the file format. This means that if a file is encoded in GB-2312 and formatted for Unix (LF), it will be considered a 'Windows mode' file; similarly, a UTF-8, Windows (CRLF) file will be considered a 'Unix mode' file.

- utf-8 bom encoded files. fconv supports utf-8 bom encoded files when reading files, but when saving files, it will only save as utf-8 (w/o bom) encoded files.

- fconv only supports GB-2312 and utf-8 encoded file conversion. So, if a file is utf-8 encoded and contains characters other than ASCII characters and Chinese simplified characters, fconv will usually refuse to convert it. 

- Since Windows mode, the file generally created under CP936 (GB-2312 encoding), which does not support Emoji. So if a UTF-8 encoded file contains Emoji characters, fconv will also refuse to conversion it.

- Some boundary rules: If file doesn't contain multi-bytes characters at all (e.g plain ASCII), fconv will think it as UTF-8 encoding. If file doesn't contain LF, fconv will think it as UNIX format.

## Credits

See go.mod file.

Use of these libraries is subject to their respective licenses.
