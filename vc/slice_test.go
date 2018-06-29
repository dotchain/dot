// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package vc_test

import (
	"encoding/json"
	"fmt"
	"github.com/dotchain/dot/vc"
)

func ExampleSliceSpliceSync_one() {
	initial := []interface{}{1, 2, 3}
	slice := vc.Slice{Version: vc.New(initial), Value: initial}
	sliced := slice.Slice(1, 2)
	fmt.Println("Sliced", sliced.Value)
	spliced := sliced.SpliceSync(1, 0, []interface{}{2.5})
	fmt.Println("Spliced", spliced.Value)
	latest, _ := slice.Version.Latest()
	json, _ := json.Marshal(latest)
	fmt.Println("Latest", string(json))

	// Output:
	// Sliced [2]
	// Spliced [2 2.5]
	// Latest [1,2,2.5,3]
}
