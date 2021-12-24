package expression

import (
	errs "gogrib/errors"
	"os"
)

var parse_file string = string("")
var top int = 0

const MAXINCLUDE int = 10

type GribArguments struct {
	expression GribExpression
	value      [80]byte
	next       *GribArguments
}

type GribAction struct {
	name         string
	op           string
	name_space   string
	next         *GribAction
	flags        uint64
	defaultKey   string
	defaultValue *GribArguments
	set          string
	debug_info   string
}
type GribActionFile struct {
	filename string
	root     *GribAction
	next     *GribActionFile
}

type GribActionFileList struct {
	first GribActionFile
	last  GribActionFile
}

func findActionFile(afl GribActionFileList, fpath string) GribActionFile {
	var act = afl.first
	for act.root != nil {
		if act.filename == fpath {
			return act
		}
		act = *act.next
	}
	return GribActionFile{}
}

func parser_include(s string) errs.Error {
	if top >= MAXINCLUDE {
		return errs.NewContextError(errs.WithMessage("top exceeds MAXINCLUDE"))
	}
	if len(s) == 0 {
		return nil
	}
	if len(parse_file) == 0 {
		parse_file = s
	} else {
		if s == "/" {
			return errs.NewContextError(errs.WithMessage("bad file to parse"))
		}
		f, err := os.OpenFile(s, os.O_RDONLY, 0666)
		if err != nil {
			return errs.NewContextError(errs.WithMessage("couldn't open file"))
		}

	}
	return nil
}

func parse(s string) errs.Error {
	parse_file = string("")
	top = 0

	err := parserInclude(s)
	if err != nil {
		return errs.NewInavalidGribError(errs.WithMessage("Grib file not found"))
	}
	err := yyparse()
	return err
}
func parseStream(s string) (*GribAction, errs.Error) {
	err := parse(s)
	if err == nil {
		return actionCreateNoop(s), nil
	}
	return nil, err

}
func gribPushActionFile(reader []GribActionFileList, af GribActionFile) {
	return
}
func GribParseFile(reader []GribActionFileList, fpath string) (GribAction, errs.Error) {
	var af GribActionFile
	if reader != nil {
		af = findActionFile(reader, fpath)
	} else {
		reader = make([]GribActionFileList, 0)
	}
	if af.root == nil {
		a, err := parseStream(fpath)
		if err != nil {
			return GribAction{}, err
		}
		af = GribActionFile{root: a, filename: fpath}
		gribPushActionFile(reader, af)
	}
	return *af.root, nil
}
