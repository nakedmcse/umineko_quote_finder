package quote

import (
	"embed"
	"math/rand/v2"
	"strings"
	"sync"

	"github.com/sahilm/fuzzy"
)

const audioDir = "internal/quote/data/audio"

//go:embed data/*.txt
var dataFS embed.FS

type Service interface {
	Search(query string, lang string, limit int, offset int, characterID string, episode int, forceFuzzy bool) SearchResponse
	GetByCharacter(lang string, characterID string, limit int, offset int, episode int) CharacterResponse
	GetByAudioID(lang string, audioID string) *ParsedQuote
	Random(lang string, characterID string, episode int) *ParsedQuote
	GetCharacters() map[string]string
	AudioFilePath(characterId string, audioId string) string
}

type service struct {
	quotes     map[string][]ParsedQuote
	quoteTexts map[string][]string
	indexer    Indexer
}

type langParseResult struct {
	lang   string
	parsed []ParsedQuote
	texts  []string
}

func NewService() Service {
	p := NewParser()

	langFiles := map[string]string{
		"en": "data/english.txt",
		"ja": "data/japanese.txt",
	}

	results := make(chan langParseResult, len(langFiles))
	var wg sync.WaitGroup

	for lang, path := range langFiles {
		wg.Add(1)
		go func(lang, path string) {
			defer wg.Done()

			data, err := dataFS.ReadFile(path)
			if err != nil {
				return
			}
			lines := strings.Split(string(data), "\n")

			parsed := p.ParseAll(lines)

			texts := make([]string, len(parsed))
			for i := 0; i < len(parsed); i++ {
				texts[i] = parsed[i].Text
			}

			results <- langParseResult{
				lang:   lang,
				parsed: parsed,
				texts:  texts,
			}
		}(lang, path)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	quotes := make(map[string][]ParsedQuote)
	texts := make(map[string][]string)

	for r := range results {
		quotes[r.lang] = r.parsed
		texts[r.lang] = r.texts
	}

	return &service{
		quotes:     quotes,
		quoteTexts: texts,
		indexer:    NewIndexer(quotes, audioDir),
	}
}

func (s *service) Search(query string, lang string, limit int, offset int, characterID string, episode int, forceFuzzy bool) SearchResponse {
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	lowerTexts := s.indexer.LowerTexts(lang)
	if quotes == nil {
		return NewSearchResponse(nil, limit, offset)
	}

	matchesFilter := func(q ParsedQuote) bool {
		if characterID != "" && q.CharacterID != characterID {
			return false
		}
		if episode > 0 && q.Episode != episode {
			return false
		}
		return true
	}

	if !forceFuzzy {
		queryLower := strings.ToLower(query)

		searchIndices := s.indexer.FilteredIndices(lang, characterID, episode)

		var exactMatches []SearchResult
		if searchIndices != nil {
			if len(searchIndices) > 5000 {
				exactMatches = concurrentExactSearch(searchIndices, lowerTexts, quotes, queryLower, matchesFilter)
			} else {
				for _, idx := range searchIndices {
					if strings.Contains(lowerTexts[idx], queryLower) {
						if matchesFilter(quotes[idx]) {
							exactMatches = append(exactMatches, NewSearchResult(quotes[idx], 100))
						}
					}
				}
			}
		} else {
			allIndices := make([]int, len(quotes))
			for i := range allIndices {
				allIndices[i] = i
			}
			exactMatches = concurrentExactSearch(allIndices, lowerTexts, quotes, queryLower, matchesFilter)
		}

		if len(exactMatches) > 0 {
			return NewSearchResponse(exactMatches, limit, offset)
		}
	}

	quoteTexts := s.quoteTexts[lang]
	matches := fuzzy.Find(query, quoteTexts)
	if len(matches) == 0 {
		return NewSearchResponse(nil, limit, offset)
	}

	topScore := matches[0].Score
	relativeThreshold := topScore / 10
	minFuzzyScore := len(query) * 100

	var fuzzyResults []SearchResult
	for i := 0; i < len(matches); i++ {
		if matches[i].Score >= relativeThreshold && matches[i].Score >= minFuzzyScore {
			if matchesFilter(quotes[matches[i].Index]) {
				fuzzyResults = append(fuzzyResults, NewSearchResult(quotes[matches[i].Index], matches[i].Score))
			}
		}
	}

	return NewSearchResponse(fuzzyResults, limit, offset)
}

func (s *service) GetByCharacter(lang string, characterID string, limit int, offset int, episode int) CharacterResponse {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil {
		return NewCharacterResponse(characterID, nil, limit, offset)
	}

	indices := s.indexer.CharacterIndices(lang, characterID)
	if len(indices) == 0 {
		return NewCharacterResponse(characterID, nil, limit, offset)
	}

	var all []ParsedQuote
	if episode <= 0 {
		all = make([]ParsedQuote, len(indices))
		for i, idx := range indices {
			all[i] = quotes[idx]
		}
	} else {
		for _, idx := range indices {
			if quotes[idx].Episode == episode {
				all = append(all, quotes[idx])
			}
		}
	}

	return NewCharacterResponse(characterID, all, limit, offset)
}

func (s *service) Random(lang string, characterID string, episode int) *ParsedQuote {
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil || len(quotes) == 0 {
		return nil
	}

	if characterID == "" && episode <= 0 {
		indices := s.indexer.NonNarratorIndices(lang)
		if len(indices) == 0 {
			return nil
		}
		pick := indices[rand.IntN(len(indices))]
		return &quotes[pick]
	}

	indices := s.indexer.FilteredIndices(lang, characterID, episode)
	if indices != nil {
		if len(indices) == 0 {
			return nil
		}
		pick := indices[rand.IntN(len(indices))]
		q := quotes[pick]
		return &q
	}

	var filtered []ParsedQuote
	for i := 0; i < len(quotes); i++ {
		if characterID != "" && quotes[i].CharacterID != characterID {
			continue
		}
		if episode > 0 && quotes[i].Episode != episode {
			continue
		}
		filtered = append(filtered, quotes[i])
	}

	if len(filtered) == 0 {
		return nil
	}

	pick := rand.IntN(len(filtered))
	return &filtered[pick]
}

func (s *service) GetByAudioID(lang string, audioID string) *ParsedQuote {
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil {
		return nil
	}

	for i := range quotes {
		if quotes[i].AudioID == audioID || strings.Contains(quotes[i].AudioID, audioID) {
			return &quotes[i]
		}
	}
	return nil
}

func (s *service) GetCharacters() map[string]string {
	return CharacterNames.GetAllCharacters()
}

func (s *service) AudioFilePath(characterId string, audioId string) string {
	return s.indexer.AudioFilePath(characterId, audioId)
}
