package tokenizer

import (
	"testing"
)

func TestSegmenterCache(t *testing.T) {
	seg1, err := getSegoSegmenter("test_dict1.txt")
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	seg2, err := getSegoSegmenter("test_dict2.txt")
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	if seg2 == seg1 {
		t.Errorf("expected to be different segmenters")
	}

	seg3, err := getSegoSegmenter("test_dict1.txt")
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	if seg3 != seg1 {
		t.Errorf("expected to be the same segmenter")
	}
}
