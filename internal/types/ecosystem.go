package types

// Determine the Eco-System I am working with
type Ecosystem int

const (
	Unknown Ecosystem = iota
	NPM
	PNPM
	Bun
	Yarn
	Go
	Cargo
	Python
	Java
)

var (
	EcosystemMapped = []string{
		"Unknown",
		"npm",
		"npm",
		"npm",
		"npm",
		"Go",
		"crates.io",
		"PyPI",
		"Maven",
	}
)
