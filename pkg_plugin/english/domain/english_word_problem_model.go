//go:generate mockery --output mock --name EnglishWordSentenceProblemModel
package domain

import (
	"context"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

const EnglishWordProblemType = "english_word"

type EnglishWordSentenceProblemModel interface {
	GetProvider() string
	GetAudioID() appD.AudioID
	GetText() string
	GetTranslated() string
	GetNote() string
	GetLang2() appD.Lang2
}

type englishWordSentenceProblemModel struct {
	Provider   string
	AudioID    appD.AudioID
	Text       string
	Lang2      appD.Lang2
	Translated string
	Note       string
}

func NewEnglishWordProblemSentenceModel(audioID appD.AudioID, text string, lang2 appD.Lang2, translated, note string) (EnglishWordSentenceProblemModel, error) {
	return &englishWordSentenceProblemModel{
		AudioID:    audioID,
		Text:       text,
		Lang2:      lang2,
		Translated: translated,
		Note:       note,
	}, nil
}

func (m *englishWordSentenceProblemModel) GetProvider() string {
	return m.Provider
}

func (m *englishWordSentenceProblemModel) GetAudioID() appD.AudioID {
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

func (m *englishWordSentenceProblemModel) GetLang2() appD.Lang2 {
	return m.Lang2
}

type EnglishWordProblemModel interface {
	appD.ProblemModel
	GetAudioID() appD.AudioID
	GetText() string
	GetPos() int
	GetPhonetic() string
	GetPresentThird() string
	GetPresentParticiple() string
	GetPastTense() string
	GetPastParticiple() string
	GetLang2() appD.Lang2
	GetTranslated() string
	GetPhrases() []EnglishPhraseProblemModel
	GetSentences() []EnglishWordSentenceProblemModel
}

type englishWordProblemModel struct {
	appD.ProblemModel
	AudioID           appD.AudioID
	Text              string
	Pos               int
	Phonetic          string
	PresentThird      string
	PresentParticiple string
	PastTense         string
	PastParticiple    string
	Lang2             appD.Lang2
	Translated        string
	Phrases           []EnglishPhraseProblemModel
	Sentences         []EnglishWordSentenceProblemModel
}

func NewEnglishWordProblemModel(problemModel appD.ProblemModel, audioID appD.AudioID, text string, pos int, phonetic string, presentThird, presentParticiple, pastTense, pastParticiple string, lang2 appD.Lang2, translated string, phrases []EnglishPhraseProblemModel, sentences []EnglishWordSentenceProblemModel) (EnglishWordProblemModel, error) {
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
		Lang2:             lang2,
		Translated:        translated,
		Phrases:           phrases,
		Sentences:         sentences,
	}, nil
}

func (m *englishWordProblemModel) GetAudioID() appD.AudioID {
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

func (m *englishWordProblemModel) GetLang2() appD.Lang2 {
	return m.Lang2
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
			"lang2":      s.GetLang2().String(),
			"note":       s.GetNote(),
		}
		sentences = append(sentences, sentence)
	}

	return map[string]interface{}{
		"text":       m.Text,
		"pos":        m.Pos,
		"lang2":      m.Lang2.String(),
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
