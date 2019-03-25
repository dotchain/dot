package myslice

// MySlice is public
type MySlice []bool
type mySlice2 []MySlice
type mySlice3 []*bool

// MySliceP is public
type MySliceP []bool
type mySlice2P []*MySliceP
type mySlice3P []*bool
