package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cherries-works/guard/internal/parser"
	"github.com/cherries-works/guard/internal/types"
	"github.com/cherries-works/guard/internal/utils"
)

const (
	SEPARATOR    = "──────────────────────────────────────────────────────"
	SUMMARY_SIZE = 48
	UNKNOWN      = "unknown"

	// MAX_ENTRIES caps how many dependencies each section lists in the
	// default output. Verbose output is never capped.
	MAX_ENTRIES = 10
)

type Analysis struct {
	Project Project

	Dependencies []Dependency
	Findings     []Finding

	DependencyCount int
}

func Analyzer(pwd string) Analysis {
	project := GetProject(pwd)

	analysis := Analysis{
		Project: project,
	}

	for ecosystem_i, ecosystem := range project.Ecosystems {
		for manifest_i, manifest := range ecosystem.Manifests {
			dependencies, findings := GetDependencies(
				ecosystem.Ecosystem,
				manifest.Path,
			)
			analysis.Project.Ecosystems[ecosystem_i].Manifests[manifest_i].DependencyCount = len(dependencies)

			analysis.Dependencies = append(
				analysis.Dependencies,
				dependencies...,
			)

			analysis.Findings = append(
				analysis.Findings,
				findings...,
			)
		}

		for _, lock := range ecosystem.Locks {
			fmt.Printf("%s\n", lock.Path)
			x := parser.ParseLock(
				ecosystem.Ecosystem,
				lock.Path,
			)
			_ = x
		}
	}

	analysis.DependencyCount = len(analysis.Dependencies)

	return analysis
}

func PrintAnalysis(analysis Analysis, verbose bool) {
	utils.Title()
	fmt.Println()

	fmt.Println("Project")
	fmt.Println(SEPARATOR)

	fmt.Printf("  Path       %s\n", analysis.Project.Root)

	fmt.Print("  Ecosystems ")

	for i, ecosystem := range analysis.Project.Ecosystems {
		if i > 0 {
			fmt.Print(", ")
		}

		fmt.Print(types.EcosystemMapped[ecosystem.Ecosystem])
	}

	fmt.Println()
	fmt.Println()

	PrintManifests(analysis.Project)
	PrintLocks(analysis.Project)
	fmt.Println()

	vulns := PrintVulnerable(analysis, verbose)
	outs := PrintOutdated(analysis, verbose)

	fmt.Println("Results")
	fmt.Println(SEPARATOR)

	fmt.Printf("  Dependencies    %d\n", analysis.DependencyCount)
	fmt.Printf("  Vulnerable      %d\n", vulns)
	fmt.Printf("  Outdated        %d\n", outs)

	fmt.Println()

	if len(analysis.Findings) > 0 {
		fmt.Println("Guard found issues.")
	} else {
		fmt.Println("No issues found.")
	}
}

func PrintTree(analysis Analysis) {
	utils.Title()
	fmt.Println()
}

func PrintLocks(project Project) (int, int, int) {
	fmt.Println("Manifests")
	fmt.Println(SEPARATOR)
	fmt.Printf("  %-15s %-25s %s\n", "Ecosystem", "Locks", "Dependencies")

	vulnerablePackages := 0
	outdatedPackages := 0
	totalDependencies := 0

	for _, ecosystem := range project.Ecosystems {
		for _, lock := range ecosystem.Locks {
			fmt.Printf(
				"  %-15s %-25s %d\n",
				types.EcosystemMapped[ecosystem.Ecosystem],
				filepath.Base(lock.Path),
				lock.DependencyCount,
			)

			totalDependencies += len(ecosystem.Dependencies)
			for _, d := range ecosystem.Dependencies {
				if d.Outdated {
					outdatedPackages++
				}
				if d.Vulnerable {
					vulnerablePackages++
				}
			}
		}
	}

	fmt.Println(SEPARATOR)

	return totalDependencies, vulnerablePackages, outdatedPackages
}

func PrintManifests(project Project) (int, int, int) {
	fmt.Println("Locks")
	fmt.Println(SEPARATOR)
	fmt.Printf("  %-15s %-25s %s\n", "Ecosystem", "Manifests", "Dependencies")

	vulnerablePackages := 0
	outdatedPackages := 0
	totalDependencies := 0

	for _, ecosystem := range project.Ecosystems {
		for _, manifest := range ecosystem.Manifests {
			fmt.Printf(
				"  %-15s %-25s %d\n",
				types.EcosystemMapped[ecosystem.Ecosystem],
				filepath.Base(manifest.Path),
				manifest.DependencyCount,
			)

			totalDependencies += len(ecosystem.Dependencies)
			for _, d := range ecosystem.Dependencies {
				if d.Outdated {
					outdatedPackages++
				}
				if d.Vulnerable {
					vulnerablePackages++
				}
			}
		}
	}

	fmt.Println(SEPARATOR)

	return totalDependencies, vulnerablePackages, outdatedPackages
}

func findFindingsById(findings []Finding, id string) []Finding {
	_findings := []Finding{}
	for _, finding := range findings {
		if finding.DependencyID == id {
			_findings = append(_findings, finding)
		}
	}
	return _findings
}

// PrintVulnerable lists every dependency with at least one known
// vulnerability. The default output keeps one line per dependency and is
// capped; verbose output adds every advisory and lists all of them.
func PrintVulnerable(analysis Analysis, verbose bool) int {
	if analysis.DependencyCount == 0 {
		return 0
	}

	fmt.Println("Vulnerable")
	fmt.Println(SEPARATOR)
	fmt.Printf("  %-30s %-15s %-15s %s\n", "Package", "Current", "Fixed in", "Advisiories")

	shown := Limit(analysis.DependencyCount, verbose)
	i := 0
	for _, dependency := range analysis.Dependencies {
		if !dependency.Vulnerable {
			continue
		}

		id := dependency.ID
		findings := findFindingsById(analysis.Findings, id)

		i++
		fmt.Printf(
			"  %-30s %-15s %-15s %d\n",
			dependency.Name,
			VersionOrUnknown(dependency.Version),
			VersionOrUnknown(findings[0].FixedVersion),
			len(findings),
		)

		if !verbose {
			continue
		}

		for _, finding := range findings {
			summary := finding.Summary
			if len(summary) > 30 {
				summary = summary[:30] + "..."
			}

			fmt.Printf(
				"    %-24s %-10s %s\n",
				finding.ID,
				finding.Severity,
				summary,
			)
		}

		fmt.Println()
	}

	PrintRemainder(analysis.DependencyCount - shown)

	if !verbose {
		fmt.Println()
	}

	return i
}

// PrintOutdated lists every dependency that is behind its latest
// published release.
func PrintOutdated(analysis Analysis, verbose bool) int {
	if analysis.DependencyCount == 0 {
		return 0
	}

	fmt.Println("Outdated")
	fmt.Println(SEPARATOR)
	fmt.Printf("  %-30s %-15s %s\n", "Package", "Current", "Latest")

	shown := Limit(analysis.DependencyCount, verbose)
	i := 0
	for _, dependency := range analysis.Dependencies[:shown] {
		if dependency.Version == dependency.LatestVersion {
			continue
		}

		i++
		fmt.Printf(
			"  %-30s %-15s %s\n",
			dependency.Name,
			VersionOrUnknown(dependency.Version),
			VersionOrUnknown(dependency.LatestVersion),
		)
	}

	PrintRemainder(analysis.DependencyCount - shown)

	fmt.Println()
	return i
}

// Limit reports how many entries a section should list.
func Limit(total int, verbose bool) int {
	if verbose || total <= MAX_ENTRIES {
		return total
	}

	return MAX_ENTRIES
}

// PrintRemainder reports the entries hidden by the default cap.
func PrintRemainder(remaining int) {
	if remaining <= 0 {
		return
	}

	fmt.Printf("  … and %d more (run with -v)\n", remaining)
}

func SeverityOf(vuln OSVVulnerability) string {
	if vuln.DatabaseSpecifics.Severity == "" {
		return UNKNOWN
	}

	return vuln.DatabaseSpecifics.Severity
}

func SummaryOf(vuln OSVVulnerability) string {
	summary := vuln.Summary

	if summary == "" {
		summary = strings.SplitN(vuln.Details, "\n", 2)[0]
	}

	summary = strings.TrimSpace(summary)

	if summary == "" {
		return "no summary available"
	}

	return Truncate(summary, SUMMARY_SIZE)
}

func Truncate(text string, size int) string {
	runes := []rune(text)

	if len(runes) <= size {
		return text
	}

	return string(runes[:size-3]) + "..."
}

func VersionOrUnknown(version string) string {
	if version == "" {
		return UNKNOWN
	}

	return version
}
