package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

const EnglishWordProblemType = "english_word"

type EnglishWordProblem interface {
	app.Problem
	GetAudioID() app.AudioID
	GetText() string
	GetPos() int
	GetPhonetic() string
	GetPresentThird() string
	GetPresentParticiple() string
	GetPastTense() string
	GetPastParticiple() string
	GetLang2() app.Lang2
	GetTranslated() string
	GetPhrases() []EnglishPhraseProblem
	GetSentences() []EnglishSentenceProblem
}

type englishWordProblem struct {
	app.Problem
	AudioID           app.AudioID
	Text              string
	Pos               int
	Phonetic          string
	PresentThird      string
	PresentParticiple string
	PastTense         string
	PastParticiple    string
	Lang              app.Lang2
	Translated        string
	Phrases           []EnglishPhraseProblem
	Sentences         []EnglishSentenceProblem
}

func NewEnglishWordProblem(problem app.Problem, audioID app.AudioID, text string, pos int, phonetic string, presentThird, presentParticiple, pastTense, pastParticiple string, lang app.Lang2, translated string, phrases []EnglishPhraseProblem, sentences []EnglishSentenceProblem) (EnglishWordProblem, error) {
	return &englishWordProblem{
		Problem:           problem,
		AudioID:           audioID,
		Text:              text,
		Pos:               pos,
		Phonetic:          phonetic,
		PresentThird:      presentThird,
		PresentParticiple: presentParticiple,
		PastTense:         pastTense,
		PastParticiple:    pastParticiple,
		Lang:              lang,
		Translated:        translated,
		Phrases:           phrases,
		Sentences:         sentences,
	}, nil
}

func (m *englishWordProblem) GetAudioID() app.AudioID {
	return m.AudioID
}

func (m *englishWordProblem) GetText() string {
	return m.Text
}

func (m *englishWordProblem) GetPos() int {
	return m.Pos
}

func (m *englishWordProblem) GetPhonetic() string {
	return m.Phonetic
}

func (m *englishWordProblem) GetPresentThird() string {
	return m.PresentThird
}

func (m *englishWordProblem) GetPresentParticiple() string {
	return m.PresentParticiple
}

func (m *englishWordProblem) GetPastTense() string {
	return m.PastTense
}

func (m *englishWordProblem) GetPastParticiple() string {
	return m.PastParticiple
}

func (m *englishWordProblem) GetLang2() app.Lang2 {
	return m.Lang
}

func (m *englishWordProblem) GetTranslated() string {
	return m.Translated
}

func (m *englishWordProblem) GetPhrases() []EnglishPhraseProblem {
	return m.Phrases
}

func (m *englishWordProblem) GetSentences() []EnglishSentenceProblem {
	return m.Sentences
}

func (m *englishWordProblem) GetProperties(cxt context.Context) map[string]interface{} {
	// fmt.Printf("m.sentences: %v\n", m.sentences[0])

	// v, _ := json.Marshal(m.sentences[0])
	// fmt.Printf("m.sentences: %v\n", string(v))

	return map[string]interface{}{
		"text":       m.Text,
		"pos":        m.Pos,
		"lang":       m.Lang.String(),
		"translated": m.Translated,
		"audioId":    m.AudioID,
		"sentences":  m.Sentences,
	}
}

// type EnglishWordProblemWithSentences interface {
// 	EnglishWordProblem
// }

// type englishWordProblemWithSentences struct {
// 	EnglishWordProblem
// 	sentences   []EnglishSentence
// 	mySentences []EnglishSentence
// }

// type ProblemAlreadyExistsError struct {
// 	message string
// }

// func NewProblemAlreadyExistsError(message string) *ProblemAlreadyExistsError {
// 	return &ProblemAlreadyExistsError{
// 		message: message,
// 	}
// }

// func (e *ProblemAlreadyExistsError) Error() string {
// 	return e.message
// }

// type ProblemNotFoundError struct {
// 	message string
// }

// func NewProblemNotFoundError(message string) *ProblemNotFoundError {
// 	return &ProblemNotFoundError{
// 		message: message,
// 	}
// }

// func (e *ProblemNotFoundError) Error() string {
// 	return e.message
// }
