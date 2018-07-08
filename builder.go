// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package dot

// BuildImage applies the sequence of operations to an initial model
// Image and returns the result.  If the old model is empty, the
// initial model is guessed from the first operation/change
//
// Invalid operations are silently ignored.
func (u Utils) BuildImage(old *ModelImage, ops []Operation) *ModelImage {
	m := &ModelImage{}

	if old != nil {
		m.Model = old.Model
		m.BasisID = old.BasisID
	}

	for _, op := range ops {
		for _, ch := range op.Changes {
			if m.Model == nil {
				if len(ch.Path) == 0 && ch.Splice != nil && ch.Splice.After != nil {
					data, ok := u.C.TryGet(ch.Splice.After)
					if !ok {
						continue
					}
					m.Model = data.Slice(0, 0)
				} else if len(ch.Path) == 0 && ch.Set != nil {
					m.Model = map[string]interface{}{}
				} else {
					continue
				}
			}
			changes := []Change{ch}
			if model, ok := u.TryApply(m.Model, changes); ok {
				m.Model = model
			}
		}
		m.BasisID = op.ID
	}
	return m
}
