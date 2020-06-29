# RFlush KNN
K-nearest neighbors search (KNN) for RFlush

### Install
````
go get github.com/anjulapaulus/rflush_knn
````

### Implementation
````
bbox := rflush.PointToBBox([2]float64{6.967660, 79.872217},1)
	//fmt.Println(bbox.Min, bbox.Max)

indexer := rflush_knn.Wrap(tr)
indexer.Nearby(rflush_knn.Box(bbox.Min, bbox.Max, false, nil),
	func(min, max [2]float64, reference string, value interface{}, dist float64) bool {
		dista := rflush_knn.Distance([2]float64{6.971322, 79.874468},min,min[0])
		if dista >= 300 {
			fmt.Println(reference,dista)
		}
		return true
	},)

````
