//go:generate mockery --output mock --name EnglishWordSentenceProblemModel
package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

const EnglishWordProblemType = "english_word"

type EnglishWordSentenceProblemModel interface {
	GetProvider() string
	GetAudioID() app.AudioID
	GetText() string
	GetTranslated() string
	GetNote() string
	GetLang2() app.Lang2
}

type englishWordSentenceProblemModel struct {
	Provider   string
	AudioID    app.AudioID
	Text       string
	Lang       app.Lang2
	Translated string
	Note       string
}

func NewEnglishWordProblemSentenceModel(audioID app.AudioID, text string, lang app.Lang2, translated, note string) (EnglishWordSentenceProblemModel, error) {
	return &englishWordSentenceProblemModel{
		AudioID:    audioID,
		Text:       text,
		Lang:       lang,
		Translated: translated,
		Note:       note,
	}, nil
}

func (m *englishWordSentenceProblemModel) GetProvider() string {
	return m.Provider
}

func (m *englishWordSentenceProblemModel) GetAudioID() app.AudioID {
	return m.AudioID
}

func (m *englishWordSentenceProblemModel) GetText() string {
	return m.Text
}

func (m *englishWordSentenceProblemModel) GetTranslated() string {
	return m.Translated
}

func (m *englishWordSentenceProblemModel) GetNote() string {
	return m.Note
}

func (m *englishWordSentenceProblemModel) GetLang2() app.Lang2 {
	return m.Lang
}

type EnglishWordProblemModel interface {
	app.ProblemModel
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
	GetPhrases() []EnglishPhraseProblemModel
	GetSentences() []EnglishWordSentenceProblemModel
}

type englishWordProblemModel struct {
	app.ProblemModel
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
	Phrases           []EnglishPhraseProblemModel
	Sentences         []EnglishWordSentenceProblemModel
}

func NewEnglishWordProblemModel(problemModel app.ProblemModel, audioID app.AudioID, text string, pos int, phonetic string, presentThird, presentParticiple, pastTense, pastParticiple string, lang app.Lang2, translated string, phrases []EnglishPhraseProblemModel, sentences []EnglishWordSentenceProblemModel) (EnglishWordProblemModel, error) {
	return &englishWordProblemModel{
		ProblemModel:      problemModel,
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

func (m *englishWordProblemModel) GetAudioID() app.AudioID {
	return m.AudioID
}

func (m *englishWordProblemModel) GetText() string {
	return m.Text
}

func (m *englishWordProblemModel) GetPos() int {
	return m.Pos
}

func (m *englishWordProblemModel) GetPhonetic() string {
	return m.Phonetic
}

func (m *englishWordProblemModel) GetPresentThird() string {
	return m.PresentThird
}

func (m *englishWordProblemModel) GetPresentParticiple() string {
	return m.PresentParticiple
}

func (m *englishWordProblemModel) GetPastTense() string {
	return m.PastTense
}

func (m *englishWordProblemModel) GetPastParticiple() string {
	return m.PastParticiple
}

func (m *englishWordProblemModel) GetLang2() app.Lang2 {
	return m.Lang
}

func (m *englishWordProblemModel) GetTranslated() string {
	return m.Translated
}

func (m *englishWordProblemModel) GetPhrases() []EnglishPhraseProblemModel {
	return m.Phrases
}

func (m *englishWordProblemModel) GetSentences() []EnglishWordSentenceProblemModel {
	return m.Sentences
}

func (m *englishWordProblemModel) GetProperties(cxt context.Context) map[string]interface{} {
	// fmt.Printf("m.sentences: %v\n", m.sentences[0])

	// v, _ := json.Marshal(m.sentences[0])
	// fmt.Printf("m.sentences: %v\n", string(v))
	sentences := make([]map[string]interface{}, 0)
	for _, s := range m.Sentences {
		sentence := map[string]interface{}{
			"text":       s.GetText(),
			"translated": s.GetTranslated(),
			"lang":       s.GetLang2().String(),
			"note":       s.GetNote(),
		}
		sentences = append(sentences, sentence)
	}

	return map[string]interface{}{
		"text":       m.Text,
		"pos":        m.Pos,
		"lang":       m.Lang.String(),
		"translated": m.Translated,
		"audioId":    m.AudioID,
		"sentences":  sentences,
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
