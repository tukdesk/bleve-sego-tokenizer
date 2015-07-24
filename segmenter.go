package tokenizer

import (
	"sync"

	"github.com/tukdesk/sego"
)

var (
	dictSegmenterMap      = map[string]*sego.Segmenter{}
	dictSegmenterMapMutex sync.Mutex
)

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
