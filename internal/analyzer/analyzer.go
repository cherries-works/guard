package analyzer

import (
	"fmt"
	"path/filepath"

	"github.com/cherries-works/guard/internal/types"
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

func PrintAnalysis(analysis Analysis) {
	fmt.Println("cherries.works Guard v0.1.0")
	fmt.Println()

	fmt.Println("Project")
	fmt.Println("──────────────────────────────────────────────────────")

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

	totalDependencies, vulnerablePackages, outdatedPackages := PrintManifests(analysis.Project)
	fmt.Println()

	fmt.Println("Results")
	fmt.Println("──────────────────────────────────────────────────────")

	fmt.Printf("  Dependencies    %d\n", totalDependencies)
	fmt.Printf("  Vulnerable      %d\n", vulnerablePackages)
	fmt.Printf("  Outdated        %d\n", outdatedPackages)

	fmt.Println()

	if vulnerablePackages > 0 || outdatedPackages > 0 {
		fmt.Println("Guard found issues.")
	} else {
		fmt.Println("No issues found.")
	}
}

func PrintManifests(project Project) (int, int, int) {
	fmt.Println("Manifests")
	fmt.Println("──────────────────────────────────────────────────────")
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

	fmt.Println("──────────────────────────────────────────────────────")

	return totalDependencies, vulnerablePackages, outdatedPackages
}
