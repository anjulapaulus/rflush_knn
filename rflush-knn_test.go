package rflush_knn

import (
	"github.com/anjulapaulus/rflush"
	"math"
	"math/rand"
	"testing"
)

func TestWrap(t *testing.T) {
	var rflush rflush.RTree

	index := Wrap(rflush)

	if index == nil {
		t.Error("Wrap function failed.")
	}
}

type tBox struct {
	min [2]float64
	max [2]float64
}

func randBoxes(N int) []tBox {
	boxes := make([]tBox, N)
	for i := 0; i < N; i++ {
		boxes[i].min[0] = rand.Float64()*360 - 180
		boxes[i].min[1] = rand.Float64()*180 - 90
		boxes[i].max[0] = boxes[i].min[0] + rand.Float64()
		boxes[i].max[1] = boxes[i].min[1] + rand.Float64()
		if boxes[i].max[0] > 180 || boxes[i].max[1] > 90 {
			i--
		}
	}
	return boxes
}

func TestDistance(t *testing.T) {
	dh := []struct {
		min, max [2]float64
		dist     float64
	}{
		{
			min:  [2]float64{6.884204, 79.892548},
			max:  [2]float64{6.884204, 79.892548},
			dist: 1060.804448838962,
		},
		{
			min:  [2]float64{6.887431, 79.887346},
			max:  [2]float64{6.887431, 79.887346},
			dist: 409.41685802572397,
		},
		{
			min:  [2]float64{6.880526, 79.882527},
			max:  [2]float64{6.880526, 79.882527},
			dist: 816.5750604553925,
		},
	}

	for _,item := range dh{
		dist := Distance([2]float64{6.887826, 79.883665},item.min,item.min[0])
		if dist != item.dist{
			t.Error("Distance Function Failed")
		}
	}
}

func TestBoxDistCalc(t *testing.T) {
	dh := []struct {
		min, max [2]float64
		dist     float64
	}{
		{
			min:  [2]float64{6.884204, 79.892548},
			max:  [2]float64{6.884204, 79.892548},
			dist:  0.00009,
		},
		{
			min:  [2]float64{6.887431, 79.887346},
			max:  [2]float64{6.887431, 79.887346},
			dist: 0.00001,
		},
		{
			min:  [2]float64{6.880526, 79.882527},
			max:  [2]float64{6.880526, 79.882527},
			dist: 0.00005,
		},
	}

	for _,item := range dh{
		dist := BoxDistCalc([2]float64{6.887826, 79.883665},[2]float64{6.887826, 79.883665},item.min,item.max,false)
		dista := math.Round(dist*100000)/100000
		if  dista!= item.dist{
			t.Error("Distance Function Failed",dista)
		}
	}
}


func TestIndex_Nearby(t *testing.T) {
	var tr rflush.RTree

	tbox:= randBoxes(1000)
	//bbox := rflush.PointToBBox([2]float64{6.967660, 79.872217},0.1)

	indexer := Wrap(tr)

	for k:=0; k<1000;k++{
		indexer.rtree.Insert(tbox[k].min, tbox[k].max,tbox[k])
	}

	var distBox []float64
		indexer.Nearby(Box(tbox[100].min, tbox[100].max, false, nil),
			func(min, max [2]float64, value interface{}, dist float64) bool {
				dista := Distance([2]float64{6.971322, 79.874468}, min, min[0])
				if dista >= 5000{
					distBox = append(distBox, dista)
				}
				return true
			}, )

	if len(distBox) != 1000{
		t.Error("Nearby Function Failed", len(distBox))
	}

}


func BenchmarkIndex_Nearby(b *testing.B) {
	var tr rflush.RTree

	tbox:= randBoxes(b.N)
	//bbox := rflush.PointToBBox([2]float64{6.967660, 79.872217},0.1)

	indexer := Wrap(tr)

	//for k:=0; k<b.N;k++{
	//	indexer.rtree.Insert(tbox[k].min, tbox[k].max,tbox[k])
	//}
	b.ResetTimer()
	//var distBox []float64

	for i:=0; i< b.N; i++ {
		indexer.Nearby(Box(tbox[i].min, tbox[i].max, false, nil),
			func(min, max [2]float64, value interface{}, dist float64) bool {
				//dista := Distance([2]float64{6.971322, 79.874468}, min, min[0])
				//if dista >= 5000 {
				//	distBox = append(distBox, dista)
				//}

				return true
			}, )
	}

}

func BenchmarkDistance(b *testing.B) {
	tbox:= randBoxes(b.N)

	for i:=0; i<b.N; i++ {
		Distance([2]float64{6.887826, 79.883665},tbox[i].min,6.887826)
	}
}

func BenchmarkBoxDistCalc(b *testing.B) {
	tbox:= randBoxes(b.N)

	testBBox := rflush.PointToBBox([2]float64{6.887826, 79.883665},0.1)
	for i:=0; i<b.N; i++ {
		BoxDistCalc(testBBox.Min,testBBox.Max,tbox[i].min,tbox[i].max,false)

	}
}
