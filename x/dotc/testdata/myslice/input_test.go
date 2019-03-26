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
	return nil
}

func valuesFormySlice2PStream() []*mySlice2P {
	return nil
}

func valuesFormySlice3Stream() []mySlice3 {
	return nil
}

func valuesFormySlice3PStream() []*mySlice3P {
	return nil
}
