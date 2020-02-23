package index

import (
	"reflect"
	"strings"
	"testing"

	"sandbox/indexer/document"
)

func TestSearch(t *testing.T) {

	type searchresult struct {
		document  string
		relevance float64
	}
	tests := []struct {
		docs     map[string]string
		search   []string
		expected []searchresult
	}{
		{
			docs: map[string]string{
				"doc:foo": "foo bar baz",
				"doc:moo": "moo mar maz",
				"doc:bar": "boo bar baz",
			},
			search: []string{"foo"},
			expected: []searchresult{
				{document: "doc:foo", relevance: 1.0},
			},
		},
		{
			docs: map[string]string{
				"doc:foo": "foo bar baz",
				"doc:moo": "moo mar maz",
				"doc:bar": "boo bar baz",
			},
			search: []string{"foo", "bar"},
			expected: []searchresult{
				{document: "doc:foo", relevance: 1.5},
				{document: "doc:bar", relevance: 0.5},
			},
		},
		{
			docs: map[string]string{
				"doc:foo": "foo bar baz",
				"doc:moo": "moo mar maz",
				"doc:bar": "boo bar baz",
			},
			search: []string{"foo", "bar", "abracadabra"},
			expected: []searchresult{
				{document: "doc:foo", relevance: 1.5},
				{document: "doc:bar", relevance: 0.5},
			},
		},
		{
			docs: map[string]string{
				"doc:foo": "foo bar baz",
				"doc:moo": "moo mar maz",
				"doc:bar": "boo bar baz",
			},
			search:   []string{"abracadabra"},
			expected: []searchresult{},
		},
	}

	for _, tt := range tests {
		ix := New()
		for name, body := range tt.docs {
			doc := document.New(name, strings.NewReader(body))
			if err := ix.Process(doc); err != nil {
				t.Fatalf("failed to process document: %q", err)
			}
		}
		res := ix.Search(tt.search)
		convres := make([]searchresult, 0, len(res))
		for _, res := range res {
			convres = append(convres, searchresult{
				document:  res.Document.Name(),
				relevance: res.Relevance,
			})
		}
		if !reflect.DeepEqual(convres, tt.expected) {
			t.Errorf("unexpected search result: got: %#v, want: %#v",
				convres, tt.expected)
		}
	}
}
