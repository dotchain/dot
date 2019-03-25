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
