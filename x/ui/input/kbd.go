// Copyright (C) 2018 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package input

// Keyboard is the interface that focus handlers should implement.
type Keyboard interface {
	Insert(ch string)
	Remove()
	ArrowRight()
	ArrowLeft()
	ShiftArrowRight()
	ShiftArrowLeft()
}
