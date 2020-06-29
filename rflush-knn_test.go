package rflush_knn

import (
	"github.com/anjulapaulus/rflush"
	"testing"
)

func TestWrap(t *testing.T) {
	var rflush rflush.RTree

	index := Wrap(rflush)

	if index == nil{
		t.Error("Wrap function failed.")
	}
}


