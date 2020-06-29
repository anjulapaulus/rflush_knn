package rflush_knn

import (
	"github.com/anjulapaulus/rflush"
	"github.com/anjulapaulus/tinyqueue"
	"math"
)

// Values that define WGS84 ellipsoid model of the Earth
const RE = 6378.137 // equatorial radius
const FE = 1 / 298.257223563 // flattening
const E2 = FE * (2 - FE)
const RAD = math.Pi / 180
const m = RAD * RE * 1000


type Index struct {
	rtree rflush.RTree
}

// wrapper function to wrap rflush.
func Wrap(rtree rflush.RTree) *Index{
	return &Index{rtree:rtree}
}

type qnode struct {
	dist  float64
	child rflush.Child
}

func (a qnode) Compare (b tinyqueue.Item) bool{
	return a.dist < b.(qnode).dist
}

// Nearby performs a kNN-type operation on the index.
// It's expected that the caller provides its own the `algo` function, which
// is used to calculate a distance to data. The `add` function should be
// called by the caller to "return" the data item along with a distance.
// The `iter` function will return all items from the smallest dist to the
// largest dist.
func (index *Index) Nearby(
	algo func(min, max [2]float64, data interface{}, item bool) (dist float64),
	iter func(min, max [2]float64, data interface{}, dist float64) bool,
) {
	var q tinyqueue.Queue
	var parent interface{}
	var children []rflush.Child

	for {
		// gather all children for parent
		children = index.rtree.Children(parent, children[:0])
		for _, child := range children {
			q.Push(qnode{
				dist:  algo(child.Min, child.Max,child.Data, child.Item),
				child: child,
			})
		}
		for {
			item := q.Pop()
			if item == nil {
				// nothing left in queue
				return
			}
			node := item.(qnode)
			if node.child.Item {
				if !iter(node.child.Min, node.child.Max,
					node.child.Data, node.dist) {
					return
				}
			} else {
				// gather more children
				parent = node.child.Data
				break
			}
		}
	}
}


// Box performs simple box-distance algorithm on rectangles.
// When wrapX is provided, the operation does a cylinder wrapping of the X value to allow
// for anti-meridian calculations.
// When itemDist is provided (not nil), it becomes the caller's responsibility to return the box-distance.
func Box(
	targetMin, targetMax [2]float64, wrapX bool,
	itemDist func(min, max [2]float64, data interface{}) float64,
) (
	algo func(min, max [2]float64, data interface{}, item bool) (dist float64),
) {
	return func(min, max [2]float64, data interface{}, item bool) (dist float64) {
		if item && itemDist != nil {
			return itemDist(min, max, data)
		}
		return BoxDistCalc(targetMin, targetMax, min, max, wrapX)
	}
}

func mmin(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

func mmax(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

// BoxDistCalc returns the distance from rectangle A to rectangle B.
// When wrapX is provided, the operation does a cylinder wrapping of the X value to allow
// for anti-meridian calculations.
func BoxDistCalc(aMin, aMax, bMin, bMax [2]float64, wrapX bool) float64 {
	var dist float64
	var squared float64

	// X
	squared = mmax(aMin[0], bMin[0]) - mmin(aMax[0], bMax[0])
	if wrapX {
		squaredLeft := mmax(aMin[0]-360, bMin[0]) - mmin(aMax[0]-360, bMax[0])
		squaredRight := mmax(aMin[0]+360, bMin[0]) - mmin(aMax[0]+360, bMax[0])
		squared = mmin(squared, mmin(squaredLeft, squaredRight))
	}
	if squared > 0 {
		dist += squared * squared
	}

	// Y
	squared = mmax(aMin[1], bMin[1]) - mmin(aMax[1], bMax[1])
	if squared > 0 {
		dist += squared * squared
	}

	return dist
}

// Given two points of the form [longitude, latitude], returns the distance.
// distance provided in metres
func Distance(a [2]float64,b [2]float64, lat float64) float64{
	//latitude
	//cos := math.Cos(a[0] * math.Pi / 180)
	//cos2 := 2*cos*cos - 1
	//cos3 := 2*cos*cos2 - cos
	//cos4 := 2*cos*cos3 - cos2
	//cos5 := 2*cos*cos4 - cos3
	//
	//Kx := 1000 * (111.41513 * cos - 0.09455*cos3 + 0.00012*cos5)
	//Ky := 1000 * (111.13209 - 0.56605*cos2 + 0.0012*cos4)
	//dx := (a[0] - b[0]) * Kx
	//dy := (a[1] - b[1]) * Ky
	//return math.Sqrt(dx*dx + dy*dy)

	//Creates a ruler instance for very fast approximations to common geodesic measurements around a certain latitude.
	coslat := math.Cos(lat * RAD)
	w2 := 1 / (1 - E2 * (1 - coslat * coslat))
	w := math.Sqrt(w2)

	// multipliers for converting longitude and latitude degrees into distance
	kx := m * w * coslat        // based on normal radius of curvature
	ky := m * w * w2 * (1 - E2) // based on meridional radius of curvature

	dx := (a[0] - b[0]) * kx
	dy := (a[1] - b[1]) * ky

	return math.Sqrt(dx * dx + dy * dy)
}
