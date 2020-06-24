package rflush_knn

import (
	"github.com/anjulapaulus/rflush"
)

type IndexerInterface interface {
	Insert(min, max [2]float64, reference string, value interface{})

	Search (min, max [2]float64, iter func(min, max [2]float64,reference string) bool)

	All() []rflush.BBox

	Bounds() (min, max [2]float64)

	Len() int

	Remove(min, max [2]float64, reference string, data interface{})

	Children(parent interface{}, reuse []Child,
	) []Child
}
