package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/terrascope/internal/hcl"
)

var topDir string
var ignoreNames []string

var defaultIgnoreNames = []string{".terraform/"}

func newProviderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "providers COMMAND",
		Short:   "a toolbox for working with Terraform providers",
		GroupID: commandGroupIDTunnelvision,
	}

	cmd.PersistentFlags().StringVarP(&topDir, "dir", "d", ".", "the directory to search")
	cmd.PersistentFlags().StringArrayVarP(&ignoreNames, "ignore", "i", []string{}, "names to ignore. `.terraform/` is appended to this list internally.")

	cmd.AddCommand(newVersionsCmd())
	cmd.AddCommand(newCacheCmd())
	cmd.AddCommand(newWhyCmd())

	return cmd
}

func newVersionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "prints out all the versions required by lockfiles in or under the top directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			ignore := append(ignoreNames, defaultIgnoreNames...)

			lockfilenames, err := findAll(".terraform.lock.hcl", topDir, ignore)
			if err != nil {
				return err
			}

			lockfiles := make([]*hcl.Lockfile, 0)
			logrus.Infof("Found %d lockfiles", len(lockfilenames))
			for _, filename := range lockfilenames {
				lf, err := hcl.ParseLockfile(filename)
				if err != nil {
					return err
				}

				logrus.Debugf("%s (%d providers)", filename, len(lf.Providers))

				lockfiles = append(lockfiles, lf)
			}

			versions := make(map[string][]string, 0)
			versionCount := 0
			for _, lf := range lockfiles {
				for _, p := range lf.Providers {
					if versions[p.ID] == nil {
						versions[p.ID] = make([]string, 0)
					}
					if !contains(versions[p.ID], p.Version) {
						versions[p.ID] = append(versions[p.ID], p.Version)
						versionCount++
					}
				}
			}
			versionJSON, err := json.MarshalIndent(versions, "", "  ")
			if err != nil {
				return err
			}
			logrus.Infof("Found %d provider versions", versionCount)
			fmt.Println(string(versionJSON))

			return nil
		},
	}

	return cmd
}

func newCacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cache",
		Short: "identifies a small set of terraform roots in the top directory " +
			"that, when applied, will cache the full set of providers required " +
			"by any root under the top directory",
		Long: "Identifies a small[1] set of terraform roots in the top directory\n" +
			"that use the full range of provider versions present in any root\n" +
			"under the top directory.\n\n" +
			"[1]: not very optimized right now, so it's not _the smallest_ set,\n" +
			"just _a small_ set.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ignore := append(ignoreNames, defaultIgnoreNames...)

			lockfilenames, err := findAll(".terraform.lock.hcl", topDir, ignore)
			if err != nil {
				return err
			}

			lockfiles := make(map[string]*hcl.Lockfile, 0)
			logrus.Infof("Found %d lockfiles", len(lockfilenames))
			for _, filename := range lockfilenames {
				lf, err := hcl.ParseLockfile(filename)
				if err != nil {
					return err
				}

				logrus.Debugf("%s (%d providers)", filename, len(lf.Providers))

				lockfiles[filename] = lf
			}

			requiredProviders := make([]string, 0)
			for _, lf := range lockfiles {
				for _, p := range lf.Providers {
					version := fmt.Sprintf("%s@%s", p.ID, p.Version)
					if !contains(requiredProviders, version) {
						requiredProviders = append(requiredProviders, version)
					}
				}
			}
			logrus.Infof("Found %d provider versions")
			if verbose {
				for _, v := range requiredProviders {
					logrus.Debug(v)
				}
			}

			rootsToApply := make([]string, 0)
			appliedProviders := make([]string, 0)
			for len(appliedProviders) < len(requiredProviders) {
				logrus.Debugf("round %d. %d providers already applied: %v", len(rootsToApply)+1, len(appliedProviders), appliedProviders)

				missingProviders := setSubtract(requiredProviders, appliedProviders)
				logrus.Debugf("  %d missing providers: %v", len(missingProviders), missingProviders)

				maxMissingProviderCount := 0
				var maxProviderFilename string
				for filename, lf := range lockfiles {
					providers := lf.CompactProviders()
					missingProviders := setIntersect(missingProviders, providers)
					logrus.Debug("    %v would apply %d providers (%v)", filename, len(missingProviders), missingProviders)
					if len(missingProviders) > maxMissingProviderCount {
						maxMissingProviderCount = len(missingProviders)
						maxProviderFilename = filename
					}
				}

				rootsToApply = append(rootsToApply, path.Dir(maxProviderFilename))
				rootProviders := lockfiles[maxProviderFilename].CompactProviders()
				logrus.Debugf("  applying %s", maxProviderFilename)
				logrus.Debugf("  because it has %d providers that are still missing. Its full setof providers: %v", maxMissingProviderCount, rootProviders)
				for _, p := range rootProviders {
					if contains(appliedProviders, p) {
						continue
					}
					appliedProviders = append(appliedProviders, p)
				}
			}

			logrus.Infof("You can apply these providers with %d roots: %v", len(rootsToApply), rootsToApply)

			for _, root := range rootsToApply {
				fmt.Println(root)
			}

			return nil
		},
	}

	return cmd
}

func newWhyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "why PROVIDER[@VERSION]",
		Short: "prints out all the roots in or under the top directory that " +
			"require the given provider",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ignore := append(ignoreNames, defaultIgnoreNames...)

			target := args[0]
			targetParts := strings.Split(target, "@")
			targetProvider := targetParts[0]
			var targetVersion string
			if len(targetParts) > 1 {
				targetVersion = targetParts[1]
			}

			lockfilenames, err := findAll(".terraform.lock.hcl", topDir, ignore)
			if err != nil {
				return err
			}

			logrus.Infof("Found %d lockfiles", len(lockfilenames))
			matches := make([]string, 0)
			for _, filename := range lockfilenames {
				lf, err := hcl.ParseLockfile(filename)
				if err != nil {
					return err
				}

				for _, p := range lf.Providers {
					if p.ID != targetProvider {
						continue
					}
					if len(targetVersion) == 0 || p.Version == targetVersion {
						matches = append(matches, fmt.Sprintf("%s requires %s@%s", path.Dir(filename), p.ID, p.Version))
					}
				}
			}

			logrus.Infof("%d roots found with the provider %s", len(matches), target)

			for _, root := range matches {
				fmt.Println(root)
			}

			return nil
		},
	}

	return cmd
}

func findAll(target, dir string, ignoreNames []string) ([]string, error) {
	found := make([]string, 0)
	err := filepath.Walk(dir,
		func(fullpath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			for _, ignoreName := range ignoreNames {
				if strings.Contains(fullpath, ignoreName) {
					return nil
				}
			}
			if info.Name() != target {
				return nil
			}

			found = append(found, fullpath)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return found, nil
}

func contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func setSubtract[T comparable](super, sub []T) []T {
	final := make([]T, 0)
	for _, el := range super {
		if !contains(sub, el) {
			final = append(final, el)
		}
	}
	return final
}

func setIntersect[T comparable](one, two []T) []T {
	final := make([]T, 0)
	for _, el := range one {
		if contains(two, el) {
			final = append(final, el)
		}
	}
	return final
}
