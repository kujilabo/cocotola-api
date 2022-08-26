package domain

import (
	"strings"

	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
)

type WordPos int

var (
	PosAdj   WordPos = 1
	PosAdv   WordPos = 2
	PosConj  WordPos = 3
	PosDet   WordPos = 4
	PosModal WordPos = 5
	PosNoun  WordPos = 6
	PosPrep  WordPos = 7
	PosPron  WordPos = 8
	PosVerb  WordPos = 9
	PosOther WordPos = 99
)

// https://docs.microsoft.com/ja-jp/rest/api/cognitiveservices/translator/translator/dictionary-lookup
func ParsePos(v string) (WordPos, error) {
	pos := strings.ToLower(v)
	switch pos {
	case "adj":
		return PosAdj, nil
	case "adv":
		return PosAdv, nil
	case "conj":
		return PosConj, nil
	case "det":
		return PosDet, nil
	case "modal":
		return PosModal, nil
	case "noun":
		return PosNoun, nil
	case "prep":
		return PosPrep, nil
	case "pron":
		return PosPron, nil
	case "verb":
		return PosVerb, nil
	default:
		return PosOther, nil
	}
}

func NewWordPos(i int) (WordPos, error) {
	if int(PosAdj) <= i && i <= int(PosVerb) {
		return WordPos(i), nil
	}
	if i == int(PosOther) {
		return WordPos(i), nil
	}
	return WordPos(0), liberrors.Errorf("invalid word pos. %d", i)
}
