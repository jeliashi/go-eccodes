package core

import (
	"os"
	"strings"
)

const ECC_PATH_DELIMITER_CHAR rune = ':'

func isEccodeDelim(r rune) bool {
	return r == ECC_PATH_DELIMITER_CHAR || r == '\000'
}

func splitEccodeString(s string) []string {
	return strings.FieldsFunc(s, isEccodeDelim)
}

func checkReadAccess(f string) bool {
	file, err := os.OpenFile(f, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return false
		}
	}
	file.Close()
	return true
}
