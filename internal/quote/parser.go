package quote

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

type Parser interface {
	ParseAll(lines []string) []ParsedQuote
}

type textRule struct {
	pattern   *regexp.Regexp
	htmlRepl  string
	plainRepl string
}

type parser struct {
	dialogueLineRegex  *regexp.Regexp
	narratorLineRegex  *regexp.Regexp
	voiceMetaRegex     *regexp.Regexp
	bracketRegex       *regexp.Regexp
	episodeMarkerRegex *regexp.Regexp
	cleanupPatterns    []string
	textRules          []textRule
}

func NewParser() Parser {
	return &parser{
		dialogueLineRegex:  regexp.MustCompile(`^d2? \[lv`),
		voiceMetaRegex:     regexp.MustCompile(`\[lv 0\*"(\d+)"\*"(\d+)"\]`),
		narratorLineRegex:  regexp.MustCompile("^d2? `"),
		bracketRegex:       regexp.MustCompile(`\[[^\]]*\]`),
		episodeMarkerRegex: regexp.MustCompile(`^new_(?:tea|ura|episode) (\d+)\r?$`),
		cleanupPatterns: []string{
			"`[@]", "`[\\]", "`[|]", "`\"", "\"`",
			"[@]", "[\\]", "[|]",
		},
		textRules: []textRule{
			{regexp.MustCompile(`\{n\}`), "<br>", " "},
			{regexp.MustCompile(`\{c:([A-Fa-f0-9]+):([^}]+)\}`), `<span style="color:#$1">$2</span>`, "$2"},
			{regexp.MustCompile(`\{f:\d+:([^}]+)\}`), `<span class="quote-name">$1</span>`, "$1"},
			{regexp.MustCompile(`\{p:\d{2,}:([^}]+)\}`), `<span class="quote-name">$1</span>`, "$1"},
			{regexp.MustCompile(`\{p:1:([^}]+)\}?`), `<span class="red-truth">$1</span>`, "$1"},
			{regexp.MustCompile(`\{p:2:([^}]+)\}?`), `<span class="blue-truth">$1</span>`, "$1"},
			{regexp.MustCompile(`\{ruby:([^:]+):([^}]+)\}`), `<ruby>$2<rp>(</rp><rt>$1</rt><rp>)</rp></ruby>`, "$2 ($1)"},
			{regexp.MustCompile(`\{i:([^}]+)\}`), `<em>$1</em>`, "$1"},
			{regexp.MustCompile(`\{a:[^:]*:(.*)\}`), "$1", "$1"},
			{regexp.MustCompile(`\{[a-z]+:[^}]*\}`), "", ""},
		},
	}
}

type ParsedQuote struct {
	Text        string `json:"text"`
	TextHtml    string `json:"textHtml"`
	CharacterID string `json:"characterId"`
	Character   string `json:"character"`
	AudioID     string `json:"audioId"`
	Episode     int    `json:"episode"`
}

func (p *parser) parseLine(line string) *ParsedQuote {
	if !p.dialogueLineRegex.MatchString(line) {
		return nil
	}

	allMatches := p.voiceMetaRegex.FindAllStringSubmatch(line, -1)
	if len(allMatches) == 0 || len(allMatches[0]) < 3 {
		return nil
	}

	characterID := allMatches[0][1]
	firstAudioID := allMatches[0][2]
	episode := p.parseEpisodeFromAudioID(firstAudioID)

	seen := map[string]bool{}
	var audioIDs []string
	for i := 0; i < len(allMatches); i++ {
		id := allMatches[i][2]
		if !seen[id] {
			seen[id] = true
			audioIDs = append(audioIDs, id)
		}
	}

	text, textHtml := p.extractText(line)
	if text == "" {
		return nil
	}

	return &ParsedQuote{
		Text:        text,
		TextHtml:    textHtml,
		CharacterID: characterID,
		Character:   CharacterNames.GetCharacterName(characterID),
		AudioID:     strings.Join(audioIDs, ", "),
		Episode:     episode,
	}
}

func (p *parser) parseEpisodeFromAudioID(audioID string) int {
	if len(audioID) < 1 {
		return 0
	}
	ep := int(audioID[0] - '0')
	if ep >= 1 && ep <= 8 {
		return ep
	}
	return 0
}

func (p *parser) extractText(line string) (string, string) {
	text := line

	for _, pattern := range p.cleanupPatterns {
		text = strings.ReplaceAll(text, pattern, "")
	}

	text = p.bracketRegex.ReplaceAllString(text, "")
	text = strings.TrimPrefix(text, "d2 ")
	text = strings.TrimPrefix(text, "d ")
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "`\"")
	text = strings.TrimSpace(text)

	plainText := text
	textHtml := html.EscapeString(text)

	for _, rule := range p.textRules {
		textHtml = rule.pattern.ReplaceAllString(textHtml, rule.htmlRepl)
		plainText = rule.pattern.ReplaceAllString(plainText, rule.plainRepl)
	}

	return plainText, textHtml
}

func (p *parser) parseNarratorLine(line string) *ParsedQuote {
	if !p.narratorLineRegex.MatchString(line) {
		return nil
	}

	text, textHtml := p.extractText(line)
	if text == "" {
		return nil
	}

	characterID := "narrator"
	character := CharacterNames.GetCharacterName("narrator")
	audioID := ""
	episode := 0

	allMatches := p.voiceMetaRegex.FindAllStringSubmatch(line, -1)
	if len(allMatches) > 0 && len(allMatches[0]) >= 3 && strings.Contains(line, "{a:") {
		characterID = allMatches[0][1]
		character = CharacterNames.GetCharacterName(characterID)
		firstAudioID := allMatches[0][2]
		episode = p.parseEpisodeFromAudioID(firstAudioID)

		seen := map[string]bool{}
		var audioIDs []string
		for i := 0; i < len(allMatches); i++ {
			id := allMatches[i][2]
			if !seen[id] {
				seen[id] = true
				audioIDs = append(audioIDs, id)
			}
		}
		audioID = strings.Join(audioIDs, ", ")
	}

	return &ParsedQuote{
		Text:        text,
		TextHtml:    textHtml,
		CharacterID: characterID,
		Character:   character,
		AudioID:     audioID,
		Episode:     episode,
	}
}

func (p *parser) ParseAll(lines []string) []ParsedQuote {
	quotes := make([]ParsedQuote, 0, len(lines)/6)
	currentEpisode := 0

	for i := 0; i < len(lines); i++ {
		if matches := p.episodeMarkerRegex.FindStringSubmatch(lines[i]); len(matches) >= 2 {
			ep, err := strconv.Atoi(matches[1])
			if err == nil {
				currentEpisode = ep
			}
			continue
		}

		parsed := p.parseLine(lines[i])
		if parsed == nil {
			parsed = p.parseNarratorLine(lines[i])
		}
		if parsed == nil || len(parsed.Text) <= 10 {
			continue
		}
		if parsed.Episode == 0 && currentEpisode > 0 {
			parsed.Episode = currentEpisode
		}
		quotes = append(quotes, *parsed)
	}

	return quotes
}
