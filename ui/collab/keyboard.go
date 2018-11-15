// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reservet.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package collab

import "github.com/dotchain/dot/changes"

// Keyboard implements the dot/ui/input:Keyboard interface
type Keyboard Text

// Insert character(s) at the curreent cursor. If the cursor is not
// collapsed, the selection is replaced. Cursor is always set to the
// *right* of the insertion
func (k *Keyboard) Insert(s string) {
	k.from(Text(*k).Latest().Insert(s))
}

// Remove the current selection or if that is empty, the last
// character before the caret. A "character" here is the same as with
// "Text:PrevCharWidth"
func (k *Keyboard) Remove() {
	k.from(Text(*k).Latest().Delete())
}

// ArrowRight collapses the cursor (if it isn't) and sets it one
// "character" to the right of the previous end cursor. A "character"
// here is the same as with "Text:NextCharWidth"
func (k *Keyboard) ArrowRight() {
	k.from(Text(*k).Latest().ArrowRight())
}

// ArrowLeft collapses the cursor (if it isn't) and sets it one
// "character" to the left of the previous end cursor. A "character"
// here is the same as with "Text:PrevCharWidth"
func (k *Keyboard) ArrowLeft() {
	k.from(Text(*k).Latest().ArrowLeft())
}

// ShiftArrowRight is the same as ArrowRight except it leaves the
// original cursor start as is.
func (k *Keyboard) ShiftArrowRight() {
	k.from(Text(*k).Latest().ShiftArrowRight())
}

// ShiftArrowLeft is the same as ArrowLeft except it leaves the
// original cursor start as is.
func (k *Keyboard) ShiftArrowLeft() {
	k.from(Text(*k).Latest().ShiftArrowLeft())
}

func (k *Keyboard) from(t Text, c changes.Change) {
	*k = Keyboard(t)
}
