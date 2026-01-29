package quote

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Indexer interface {
	LowerTexts(lang string) []string
	FilteredIndices(lang string, characterID string, episode int) []int
	CharacterIndices(lang string, characterID string) []int
	NonNarratorIndices(lang string) []int
	AudioFilePath(characterId string, audioId string) string
}

type indexer struct {
	quoteLowerTexts  map[string][]string
	characterIndex   map[string]map[string][]int
	episodeIndex     map[string]map[int][]int
	nonNarratorIndex map[string][]int
	quotes           map[string][]ParsedQuote
	audioDir         string
}

type langIndexResult struct {
	lang           string
	lowerTexts     []string
	charIdx        map[string][]int
	epIdx          map[int][]int
	nonNarratorIdx []int
}

func NewIndexer(quotes map[string][]ParsedQuote, audioDir string) Indexer {
	results := make(chan langIndexResult, len(quotes))
	var wg sync.WaitGroup

	for lang, parsed := range quotes {
		wg.Add(1)
		go func(lang string, parsed []ParsedQuote) {
			defer wg.Done()

			lowerTexts := make([]string, len(parsed))
			charIdx := make(map[string][]int)
			epIdx := make(map[int][]int)
			var nonNarratorIdx []int

			for i := 0; i < len(parsed); i++ {
				lowerTexts[i] = strings.ToLower(parsed[i].Text)
				charIdx[parsed[i].CharacterID] = append(charIdx[parsed[i].CharacterID], i)
				if parsed[i].Episode > 0 {
					epIdx[parsed[i].Episode] = append(epIdx[parsed[i].Episode], i)
				}
				if parsed[i].CharacterID != "narrator" {
					nonNarratorIdx = append(nonNarratorIdx, i)
				}
			}

			results <- langIndexResult{
				lang:           lang,
				lowerTexts:     lowerTexts,
				charIdx:        charIdx,
				epIdx:          epIdx,
				nonNarratorIdx: nonNarratorIdx,
			}
		}(lang, parsed)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	idx := &indexer{
		quoteLowerTexts:  make(map[string][]string),
		characterIndex:   make(map[string]map[string][]int),
		episodeIndex:     make(map[string]map[int][]int),
		nonNarratorIndex: make(map[string][]int),
		quotes:           quotes,
		audioDir:         audioDir,
	}

	for r := range results {
		idx.quoteLowerTexts[r.lang] = r.lowerTexts
		idx.characterIndex[r.lang] = r.charIdx
		idx.episodeIndex[r.lang] = r.epIdx
		idx.nonNarratorIndex[r.lang] = r.nonNarratorIdx
	}

	return idx
}

func (idx *indexer) LowerTexts(lang string) []string {
	return idx.quoteLowerTexts[lang]
}

func (idx *indexer) CharacterIndices(lang string, characterID string) []int {
	langCharIdx := idx.characterIndex[lang]
	if langCharIdx == nil {
		return nil
	}
	return langCharIdx[characterID]
}

func (idx *indexer) AudioFilePath(characterId string, audioId string) string {
	if idx.audioDir == "" {
		return ""
	}
	path := filepath.Join(idx.audioDir, characterId, audioId+".ogg")
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

func (idx *indexer) NonNarratorIndices(lang string) []int {
	return idx.nonNarratorIndex[lang]
}

func (idx *indexer) FilteredIndices(lang string, characterID string, episode int) []int {
	hasChar := characterID != ""
	hasEp := episode > 0

	if !hasChar && !hasEp {
		return nil
	}

	if hasChar && !hasEp {
		langCharIdx := idx.characterIndex[lang]
		if langCharIdx == nil {
			return []int{}
		}
		return langCharIdx[characterID]
	}

	if !hasChar && hasEp {
		langEpIdx := idx.episodeIndex[lang]
		if langEpIdx == nil {
			return []int{}
		}
		return langEpIdx[episode]
	}

	langCharIdx := idx.characterIndex[lang]
	if langCharIdx == nil {
		return []int{}
	}
	charIndices := langCharIdx[characterID]
	quotes := idx.quotes[lang]

	var result []int
	for _, i := range charIndices {
		if quotes[i].Episode == episode {
			result = append(result, i)
		}
	}
	return result
}
