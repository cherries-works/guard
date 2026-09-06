package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cherries-works/guard/internal/types"
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

	Vulnerable []Dependency
	Outdated   []Dependency

	DependencyCount int
}

func Analyzer(pwd string) Analysis {
	project := GetProject(pwd)

	analysis := Analysis{
		Project: project,
	}

	for _, ecosystem := range project.Ecosystems {
		for _, manifest := range ecosystem.Manifests {
			dependencies := GetDependencies(
				ecosystem.Ecosystem,
				manifest,
			)

			analysis.Dependencies = append(
				analysis.Dependencies,
				dependencies...,
			)

			for _, dependency := range dependencies {
				if dependency.Status.Vulnerable {
					analysis.Vulnerable = append(
						analysis.Vulnerable,
						dependency,
					)
				}

				if dependency.Status.Outdated {
					analysis.Outdated = append(
						analysis.Outdated,
						dependency,
					)
				}
			}
		}
	}

	analysis.DependencyCount = len(analysis.Dependencies)

	return analysis
}

func PrintAnalysis(analysis Analysis, verbose bool) {
	fmt.Println("cherries.works Guard v0.1.0")
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
	fmt.Println()

	PrintVulnerable(analysis.Vulnerable, verbose)
	PrintOutdated(analysis.Outdated, verbose)

	fmt.Println("Results")
	fmt.Println(SEPARATOR)

	fmt.Printf("  Dependencies    %d\n", analysis.DependencyCount)
	fmt.Printf("  Vulnerable      %d\n", len(analysis.Vulnerable))
	fmt.Printf("  Outdated        %d\n", len(analysis.Outdated))

	fmt.Println()

	if len(analysis.Vulnerable) > 0 || len(analysis.Outdated) > 0 {
		fmt.Println("Guard found issues.")
	} else {
		fmt.Println("No issues found.")
	}
}

func PrintManifests(project Project) (int, int, int) {
	fmt.Println("Manifests")
	fmt.Println(SEPARATOR)
	fmt.Printf("  %-15s %-25s %s\n", "Ecosystem", "Manifest", "Dependencies")

	vulnerablePackages := 0
	outdatedPackages := 0
	totalDependencies := 0

	for _, ecosystem := range project.Ecosystems {
		for _, manifest := range ecosystem.Manifests {
			dependencies := GetDependencies(
				ecosystem.Ecosystem,
				manifest,
			)

			fmt.Printf(
				"  %-15s %-25s %d\n",
				types.EcosystemMapped[ecosystem.Ecosystem],
				filepath.Base(manifest),
				len(dependencies),
			)

			totalDependencies += len(dependencies)
			for _, d := range dependencies {
				if d.Status.Outdated {
					outdatedPackages++
				}
				if d.Status.Vulnerable {
					vulnerablePackages++
				}
			}
		}
	}

	fmt.Println(SEPARATOR)

	return totalDependencies, vulnerablePackages, outdatedPackages
}

// PrintVulnerable lists every dependency with at least one known
// vulnerability. The default output keeps one line per dependency and is
// capped; verbose output adds every advisory and lists all of them.
func PrintVulnerable(dependencies []Dependency, verbose bool) {
	if len(dependencies) == 0 {
		return
	}

	fmt.Println("Vulnerable")
	fmt.Println(SEPARATOR)
	fmt.Printf("  %-30s %-15s %-15s %s\n", "Package", "Current", "Fixed in", "Advisories")

	shown := Limit(len(dependencies), verbose)

	for _, dependency := range dependencies[:shown] {
		fmt.Printf(
			"  %-30s %-15s %-15s %d\n",
			dependency.Name,
			VersionOrUnknown(dependency.Status.CurrentVersion),
			VersionOrUnknown(FixedVersion(dependency)),
			len(dependency.Status.Vulns),
		)

		if !verbose {
			continue
		}

		for _, vuln := range dependency.Status.Vulns {
			fmt.Printf(
				"    %-24s %-10s %s\n",
				vuln.ID,
				SeverityOf(vuln),
				SummaryOf(vuln),
			)
		}

		fmt.Println()
	}

	PrintRemainder(len(dependencies) - shown)

	if !verbose {
		fmt.Println()
	}
}

// PrintOutdated lists every dependency that is behind its latest
// published release.
func PrintOutdated(dependencies []Dependency, verbose bool) {
	if len(dependencies) == 0 {
		return
	}

	fmt.Println("Outdated")
	fmt.Println(SEPARATOR)
	fmt.Printf("  %-30s %-15s %s\n", "Package", "Current", "Latest")

	shown := Limit(len(dependencies), verbose)

	for _, dependency := range dependencies[:shown] {
		fmt.Printf(
			"  %-30s %-15s %s\n",
			dependency.Name,
			VersionOrUnknown(dependency.Status.CurrentVersion),
			VersionOrUnknown(dependency.Status.LatestVersion),
		)
	}

	PrintRemainder(len(dependencies) - shown)

	fmt.Println()
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

// FixedVersion reports the first patched release advertised by the
// advisories already attached to the dependency. It falls back to the
// latest known release when no advisory declares a fix.
func FixedVersion(dependency Dependency) string {
	for _, vuln := range dependency.Status.Vulns {
		for _, affected := range vuln.Affected {
			if !strings.EqualFold(affected.Package.Name, dependency.Name) {
				continue
			}

			for _, affectedRange := range affected.Ranges {
				for _, event := range affectedRange.Events {
					if event.Fixed != "" {
						return event.Fixed
					}
				}
			}
		}
	}

	return dependency.Status.LatestVersion
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
