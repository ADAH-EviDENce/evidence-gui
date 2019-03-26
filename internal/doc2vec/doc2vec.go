// Package doc2vec implements a nearest-neighbor index based on doc2vec.
//
// Training the doc2vec model is done offline by a Python program.
package doc2vec

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/knaw-huc/evidence-gui/internal/vectors"
	"github.com/knaw-huc/evidence-gui/internal/vp"
)

// A Document is a document id (reference to Elasticsearch)
// combined with its a doc2vec vector.
type Document struct {
	id     string
	vector vectors.Normalized
}

// An Index contains doc2vec vectors for documents and allows nearest neighbor
// queries.
type Index struct {
	docs map[string]*Document

	// The actual index structure is a VP-tree using Euclidean distance.
	// Minimizing Euclidean distance is equivalent to maximizing cosine
	// similarity.
	tree *vp.Tree
}

func NewIndexFromCSV(filename string) (*Index, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byid := make(map[string]*Document)
	docs := make([]interface{}, 0)

	r := csv.NewReader(f)
loop:
	for {
		record, err := r.Read()
		switch err {
		case nil:
		case io.EOF:
			break loop
		default:
			return nil, err
		}

		vector := make([]float32, 0)
		for _, f := range record[1:] {
			x, err := strconv.ParseFloat(f, 32)
			if err != nil {
				return nil, err
			}
			vector = append(vector, float32(x))
		}

		doc := &Document{
			id:     record[0],
			vector: vectors.NewNormalized(vector),
		}
		byid[doc.id] = doc
		docs = append(docs, doc)
	}

	tree, err := vp.New(nil, distance, docs)
	if err != nil {
		return nil, err
	}

	return &Index{docs: byid, tree: tree}, nil
}

// Performs a nearest-neighbors query for the document with the given id.
// The results offset through offset+size are returned.
func (idx *Index) Nearest(ctx context.Context, id string, offset, size int, exclude []string) ([]string, error) {
	doc, ok := idx.docs[id]
	if !ok {
		return nil, fmt.Errorf("no document with id %q", id)
	}

	excludeSet := make(map[string]struct{})
	for _, id := range exclude {
		excludeSet[id] = struct{}{}
	}

	near, err := idx.tree.Search(ctx, doc, size+offset, math.Inf(+1), func(x interface{}) bool {
		_, ok := excludeSet[x.(*Document).id]
		return !ok
	})
	if err != nil {
		return nil, err
	}

	offset = min(offset, len(near))
	end := min(offset+size, len(near))

	ids := make([]string, 0, end-offset)
	for _, n := range near[offset:end] {
		ids = append(ids, n.Point.(*Document).id)
	}

	return ids, nil
}

// Euclidean distance between the vectors d and q (document and query, but
// the order is irrelevant).
func distance(d, q interface{}) float64 {
	return vectors.Distance(d.(*Document).vector, q.(*Document).vector)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}