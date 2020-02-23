package index

import (
	"bufio"
	"sandbox/indexer/document"
	"sort"
	"strings"
)

type docFreq map[string]int

type termFreq struct {
	freq int
	docs docFreq
}

type Index struct {
	terms map[string]*termFreq
	docs  map[string]*document.Document
}

func New() *Index {
	return &Index{
		terms: make(map[string]*termFreq),
		docs:  make(map[string]*document.Document),
	}
}

func (i *Index) Process(doc *document.Document) error {
	scanner := bufio.NewScanner(doc.Reader())
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		token := tokenize(scanner.Text())
		i.reverseMap(doc, token)
	}
	i.docs[doc.Name()] = doc
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

type SearchResult struct {
	Document  *document.Document
	Relevance float64
}

func (i *Index) Search(terms []string) []SearchResult {
	res := make([]SearchResult, 0)
	docrel := make(map[string]float64)
	for _, term := range terms {
		if _, ok := i.terms[term]; !ok {
			continue
		}
		idf := float64(i.terms[term].freq)
		docs := i.terms[term].docs
		for name, tf := range docs {
			docrel[name] += float64(tf) / idf
		}
	}

	for name, rel := range docrel {
		res = append(res, SearchResult{
			Document:  i.docs[name],
			Relevance: rel,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Relevance > res[j].Relevance
	})

	return res
}

func (i *Index) reverseMap(doc *document.Document, token string) {
	if _, ok := i.terms[token]; !ok {
		i.terms[token] = &termFreq{
			docs: make(map[string]int),
		}
	}
	i.terms[token].freq++

	i.terms[token].docs[doc.Name()]++
}

func tokenize(text string) string {
	text = strings.Trim(text, "[](),.\"'`")
	text = strings.ToLower(text)
	chunks := strings.Split(text, ".")
	text = chunks[0]
	return text
}
