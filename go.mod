module github.com/dotchain/dot

go 1.12

require (
	github.com/etcd-io/bbolt v1.3.2
	github.com/google/go-cmp v0.2.0
	github.com/lib/pq v1.1.0
	github.com/sergi/go-diff v1.0.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/tvastar/test v0.0.0-20190408215541-5e6ef1905826
	github.com/tvastar/toc v0.0.0-20190328211025-65f5d4ff75b4 // indirect
	golang.org/x/tools v0.0.0-20190420181800-aa740d480789
	gopkg.in/russross/blackfriday.v2 v2.0.0-00010101000000-000000000000 // indirect
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday v2.0.0+incompatible
