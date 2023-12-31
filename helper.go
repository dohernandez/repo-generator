package generator

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// Receiver returns the receiver of the repo.
func Receiver(r string) string {
	words := strings.Split(strcase.ToSnake(r), "_")

	return string(strings.ToLower(words[len(words)-1])[0])
}
