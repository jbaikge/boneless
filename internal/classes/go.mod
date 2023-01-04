module github.com/jbaikge/boneless/internal/classes

go 1.19

require (
	github.com/jbaikge/boneless/internal/common v0.0.0-00010101000000-000000000000
	github.com/zeebo/assert v1.3.1
	golang.org/x/exp v0.0.0-20221230185412-738e83a70c30
)

require github.com/rs/xid v1.4.0 // indirect

replace github.com/jbaikge/boneless/internal/common => ../common/
