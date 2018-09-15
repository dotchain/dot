// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package ops_test

import (
	"context"
	"fmt"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/x/idgen"
	"github.com/dotchain/dot/x/types"
)

func Example_sync() {
	store := MemStore(nil)
	xformed := ops.Transformed(store)

	client1, client2 := changes.NewStream(), changes.NewStream()

	s1 := ops.NewSync(xformed, -1, client1, idgen.New)
	s2 := ops.NewSync(xformed, -1, client2, idgen.New)

	defer s1.Close()
	defer s2.Close()

	client1 = client1.Append(changes.Splice{0, types.S8(""), types.S8("Hello ")})
	client2 = client2.Append(changes.Splice{0, types.S8(""), types.S8("World")})

	if err := s1.Fetch(context.Background(), 100); err != nil {
		panic(err)
	}

	c, _ := client1.Next()
	if c != (changes.Splice{6, types.S8(""), types.S8("World")}) {
		fmt.Println("Unexpected client1 change", c)
	}

	c, _ = client2.Next()
	if c != nil {
		fmt.Println("Sync fetched before calling Fetch", c)
	}

	if err := s2.Fetch(context.Background(), 100); err != nil {
		panic(err)
	}

	c, _ = client2.Next()
	if c != (changes.Splice{0, types.S8(""), types.S8("Hello ")}) {
		fmt.Println("Unexpected client2  change", c)
	}

	// Output:
}
