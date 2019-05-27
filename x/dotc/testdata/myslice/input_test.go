package myslice

func valuesForMySliceStream() []MySlice {
	return []MySlice{
		MySlice{},
		MySlice{true},
		MySlice{false},
		MySlice{true, false},
	}
}

func valuesForMySlicePStream() []*MySliceP {
	return []*MySliceP{
		&MySliceP{},
		&MySliceP{true},
		&MySliceP{false},
		&MySliceP{true, false},
	}
}

func valuesFormySlice2Stream() []mySlice2 {
	return []mySlice2{
		mySlice2{},
		mySlice2(valuesForMySliceStream()[:1]),
		mySlice2(valuesForMySliceStream()[:2]),
		mySlice2(valuesForMySliceStream()[:3]),
	}
}

func valuesFormySlice2PStream() []*mySlice2P {
	values := valuesForMySlicePStream()
	v1, v2, v3 := mySlice2P(values[:1]), mySlice2P(values[:2]), mySlice2P(values[:3])

	return []*mySlice2P{&mySlice2P{}, &v1, &v2, &v3}
}

func valuesFormySlice3Stream() []mySlice3 {
	vTrue, vFalse := true, false
	return []mySlice3{
		mySlice3{},
		mySlice3{&vTrue},
		mySlice3{&vFalse},
		mySlice3{&vTrue, &vFalse},
	}
}

func valuesFormySlice3PStream() []*mySlice3P {
	vTrue, vFalse := true, false
	values := mySlice3P([]*bool{&vTrue, &vFalse, &vTrue, &vFalse})
	v1, v2, v3, v4 := values[:0], values[:1], values[:2], values[:3]

	return []*mySlice3P{&v1, &v2, &v3, &v4}
}
