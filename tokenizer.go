package tokenizer

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"github.com/tukdesk/sego"
)

const (
	Name = "huichen/sego"
)

var (
	_ analysis.Tokenizer = &SegoTokenizer{}

	ideographRegexp = regexp.MustCompile(`\p{Han}+`)

	dictSegmenterMap      = map[string]*sego.Segmenter{}
	dictSegmenterMapMutex sync.Mutex
)

type SegoTokenizer struct {
	segmenter     *sego.Segmenter
	nested        bool
	caseSensitive bool
}

func NewSegoTokenizer(dictFiles string, nested, caseSensitive bool) (*SegoTokenizer, error) {
	segmenter, err := getSegoSegmenter(dictFiles)
	if err != nil {
		return nil, err
	}

	return &SegoTokenizer{
		segmenter:     segmenter,
		nested:        nested,
		caseSensitive: caseSensitive,
	}, nil
}

func getSegoSegmenter(dictFiles string) (*sego.Segmenter, error) {
	dictSegmenterMapMutex.Lock()
	defer dictSegmenterMapMutex.Unlock()

	if segmenter, ok := dictSegmenterMap[dictFiles]; ok {
		return segmenter, nil
	}

	segmenter := new(sego.Segmenter)
	if err := segmenter.LoadDictionary(dictFiles); err != nil {
		return nil, err
	}

	dictSegmenterMap[dictFiles] = segmenter

	return segmenter, nil
}

func (this *SegoTokenizer) Tokenize(b []byte) analysis.TokenStream {
	stream := make(analysis.TokenStream, 0)
	pos := 1

	segments := this.segmenter.SegmentAdv(b, false, this.caseSensitive)
	for _, segment := range segments {
		p := &segment
		stream, pos = appendToTokenStreams(stream, p, p.Start(), pos, this.nested, true)
	}
	return stream
}

func appendToTokenStreams(stream analysis.TokenStream, segment *sego.Segment, start, pos int, nested, top bool) (analysis.TokenStream, int) {
	if nested && len(segment.Token().Segments()) > 0 {
		for _, one := range segment.Token().Segments() {
			stream, pos = appendToTokenStreams(stream, one, start+one.Start(), pos, nested, false)
		}
	}

	if top || !isFake(segment) {
		token := &analysis.Token{
			Term:     []byte(segment.Token().Text()),
			Start:    start,
			End:      start + segment.End() - segment.Start(),
			Position: pos,
			Type:     tokenType(segment.Token().Text()),
		}
		stream = append(stream, token)
		pos++
	}

	return stream, pos
}

func isFake(segment *sego.Segment) bool {
	return segment.Token().Frequency() == 1 && segment.Token().Pos() == "x"
}

func tokenType(s string) analysis.TokenType {
	if ideographRegexp.MatchString(s) {
		return analysis.Ideographic
	}

	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return analysis.Numeric
	}

	return analysis.AlphaNumeric
}

func SegoTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	dictFiles, ok := config["files"].(string)
	if !ok {
		return nil, fmt.Errorf("dictionary file paths required")
	}

	nested, ok := config["nested"].(bool)
	if !ok {
		nested = true
	}

	caseSensitive, ok := config["case"].(bool)
	if !ok {
		caseSensitive = true
	}

	return NewSegoTokenizer(dictFiles, nested, caseSensitive)
}

func init() {
	registry.RegisterTokenizer(Name, SegoTokenizerConstructor)
}
