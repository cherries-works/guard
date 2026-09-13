package parser

import (
	"github.com/cherries-works/guard/internal/lock"
	"github.com/cherries-works/guard/internal/manifest"
	"github.com/cherries-works/guard/internal/types"
)

func ParseManifest(ecosystem types.Ecosystem, path string) []string {
	return manifest.ParseManifest(ecosystem, path)
}

func ParseLock(ecosystem types.Ecosystem, path string) []string {
	return lock.ParseLock(ecosystem, path)
}
