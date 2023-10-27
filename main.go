package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	VERSION = "1.0.0"
)

const (
	NONEMODE = iota
	AUTOMODE
	WINDOWSMODE
	UNIXMODE
)

const (
	INVALIDFORMAT = iota
	DOSFORMAT
	UNIXFORMAT
)

type Fconv struct {
	mode int
}

type Result struct {
	// 0 - from, 1 - to
	Encoding [2]string
	Format   [2]int
}

func main() {
	optPrint := flag.Bool("p", false, "Print encoding and format")
	optUnix := flag.Bool("u", false, "Unix mode")
	optWindows := flag.Bool("w", false, "Windows mode")
	optVer := flag.Bool("v", false, "Print version and exit")

	flag.Parse()

	if *optVer {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	var mode = AUTOMODE
	if *optPrint {
		mode = NONEMODE
	} else if *optUnix && *optWindows {
		fmt.Printf("Cannot use -u and -w at the same time\n")
		os.Exit(1)
	} else if *optWindows {
		mode = WINDOWSMODE
	} else if *optUnix {
		mode = UNIXMODE
	}

	f := NewFconv(mode)

	args := flag.Args()
	for _, arg := range args {
		r, err := f.ConvertFile(arg)
		if err != nil {
			fmt.Printf("%s: %v\n", arg, err)
			continue
		}

		if (r.Encoding[0] == r.Encoding[1]) && (r.Format[0] == r.Format[1]) {
			fmt.Printf("%s: %s %s\n", arg, r.Encoding[0], formatToString(r.Format[0]))
		} else {
			fmt.Printf("%s: %s %s -> %s %s\n", arg,
				r.Encoding[0], formatToString(r.Format[0]),
				r.Encoding[1], formatToString(r.Format[1]))
		}
	}
}

func NewFconv(mode int) *Fconv {
	return &Fconv{
		mode: mode,
	}
}

func (f *Fconv) ConvertFile(path string) (Result, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Result{}, err
	}

	// skip bom
	if len(content) > 3 && (content[0] == 0xef && content[1] == 0xbb && content[2] == 0xbf) {
		content = content[3:]
	}

	encoding, err := detectEncoding(content)
	if err != nil {
		return Result{}, err
	}

	format := detectFormat(content)
	if format == INVALIDFORMAT {
		return Result{}, fmt.Errorf("unknown format")
	}

	var toEncoding string
	var toFormat int

	switch f.mode {
	case NONEMODE:
		toEncoding = encoding
		toFormat = format
	case AUTOMODE:
		if encoding == "UTF-8" {
			toEncoding = "GB-18030"
			toFormat = DOSFORMAT
		} else if encoding == "GB-18030" {
			toEncoding = "UTF-8"
			toFormat = UNIXFORMAT
		}
	case WINDOWSMODE:
		toEncoding = "GB-18030"
		toFormat = DOSFORMAT
	case UNIXMODE:
		toEncoding = "UTF-8"
		toFormat = UNIXFORMAT
	}

	content, c1, err := changeEncoding(content, encoding, toEncoding)
	if err != nil {
		return Result{}, err
	}

	content, c2, err := changeFormat(content, format, toFormat)
	if err != nil {
		return Result{}, err
	}

	if c1 || c2 {
		err := os.WriteFile(path, content, 0644)
		if err != nil {
			return Result{}, err
		}
	}

	return Result{
		Encoding: [2]string{encoding, toEncoding},
		Format:   [2]int{format, toFormat},
	}, nil
}

func detectEncoding(content []byte) (string, error) {
	possible, err := chardet.NewTextDetector().DetectAll(content)
	if err != nil {
		return "", err
	}

	var encoding string
	for _, p := range possible {
		if p.Charset == "Shift_JIS" || p.Charset == "GB-18030" {
			encoding = "GB-18030"
			break
		} else if p.Charset == "UTF-8" {
			encoding = "UTF-8"
			break
		}
	}

	if encoding == "" {
		return encoding, fmt.Errorf("unknown encoding: %v", possible)
	}

	if encoding == "UTF-8" {
		for len(content) > 0 {
			r, size := utf8.DecodeRune(content)

			content = content[size:]

			if r < 0x80 {
				continue
			}

			// FIXME: not accurate
			if !unicode.Is(unicode.Han, r) {
				return "", fmt.Errorf("may contain non-chinese characters, unsafe to convert")
			}
		}
	}

	return encoding, nil
}

func detectFormat(content []byte) int {
	if bytes.Contains(content, []byte("\r\n")) {
		return DOSFORMAT
	} else if bytes.Contains(content, []byte("\n")) {
		return UNIXFORMAT
	} else {
		return INVALIDFORMAT
	}
}

func changeEncoding(content []byte, from string, to string) ([]byte, bool, error) {
	var transformer transform.Transformer
	var err error

	if from == to {
		return content, false, nil
	} else if from == "UTF-8" && to == "GB-18030" {
		transformer = simplifiedchinese.GB18030.NewEncoder()
	} else if from == "GB-18030" && to == "UTF-8" {
		transformer = simplifiedchinese.GB18030.NewDecoder()
	} else {
		return content, false, fmt.Errorf("cannot do encoding")
	}

	// FIXME: use a better way to estimate the size of output
	out := make([]byte, len(content)*2)
	n, _, err := transformer.Transform(out, content, true)
	if err != nil {
		return nil, false, err
	}

	return out[:n], true, nil
}

func changeFormat(content []byte, from int, to int) ([]byte, bool, error) {
	if from == to {
		return content, false, nil
	}

	switch from {
	case DOSFORMAT:
		if to == UNIXFORMAT {
			return bytes.Replace(content, []byte("\r\n"), []byte("\n"), -1), true, nil
		}
	case UNIXFORMAT:
		if to == DOSFORMAT {
			return bytes.Replace(content, []byte("\n"), []byte("\r\n"), -1), true, nil
		}
	}

	return nil, false, fmt.Errorf("invalid format")
}

func formatToString(format int) string {
	switch format {
	case DOSFORMAT:
		return "DOS"
	case UNIXFORMAT:
		return "UNIX"
	default:
		return "UNKNOWN"
	}
}
