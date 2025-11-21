package parquet

import (
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/parquet", new(RootModule))
}

// RootModule is the global module instance that will create module
// instances for each VU.
type RootModule struct{}

// Parquet represents an instance of the module for every VU.
type Parquet struct {
	vu    modules.VU
	cache *ReaderCache
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &Parquet{}
	_ modules.Module   = &RootModule{}
)

// NewModuleInstance implements the modules.Module interface and returns
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Parquet{
		vu:    vu,
		cache: NewReaderCache(),
	}
}

// Exports implements the modules.Instance interface and returns
// the exports of the JS module.
func (p *Parquet) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"read":        p.Read,
			"readChunked": p.ReadChunked,
			"getSchema":   p.GetSchema,
			"getMetadata": p.GetMetadata,
			"close":       p.Close,
		},
	}
}
