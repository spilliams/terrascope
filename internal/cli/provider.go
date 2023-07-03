package cli

import (
	"crypto/sha256"
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
		Use:   "provider COMMAND",
		Short: "A toolbox for working with Terraform providers",
	}

	cmd.PersistentFlags().StringVar(&topDir, "dir", ".", "the directory to search")
	cmd.PersistentFlags().StringArrayVarP(&ignoreNames, "ignore", "i", []string{}, "names to ignore. `.terraform/` is appended to this list internally.")

	cmd.AddCommand(newProviderCacheCmd())
	cmd.AddCommand(newProviderHashesCmd())
	cmd.AddCommand(newProviderVersionsCmd())
	cmd.AddCommand(newProviderWhyCmd())

	return cmd
}

func newProviderCacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "identifies a small set of terraform roots in the top directory " + "that, when applied, will cache the full set of providers required " + "by any root under the top directory",
		Long:  "Identifies a small[1] set of terraform roots in the top directory\n" + "that use the full range of provider versions present in any root\n" + "under the top directory.\n\n" + "[1]: not very optimized right now, so it's not _the smallest_ set,\n" + "just _a small_ set.",
		RunE: func(cmd *cobra.Command, args []string) error {
			lockfiles, err := getLockfiles()
			if err != nil {
				return err
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
			logrus.Infof("Found %d required %s", len(requiredProviders), pluralize("provider", "providers", len(requiredProviders)))
			if verbose {
				for _, v := range requiredProviders {
					logrus.Debug(v)
				}
			}
			rootsToApply := make([]string, 0)
			appliedProviders := make([]string, 0)
			for len(appliedProviders) < len(requiredProviders) {
				logrus.Debugf("round %d. %d %s already applied: %v",
					len(rootsToApply)+1,
					len(appliedProviders),
					pluralize("provider", "providers", len(appliedProviders)),
					appliedProviders)
				missingProviders := setSubtract(requiredProviders, appliedProviders)
				logrus.Debugf("  %d missing %s: %v",
					len(missingProviders),
					pluralize("provider", "providers", len(missingProviders)),
					missingProviders)
				maxMissingProviderCount := 0
				var maxProviderFilename string
				for filename, lf := range lockfiles {
					providers := lf.CompactProviders()
					missingProviders := setIntersect(missingProviders, providers)
					logrus.Debugf("    %s would apply %d %s (%v)",
						filename,
						len(missingProviders),
						pluralize("provider", "providers", len(missingProviders)),
						missingProviders)
					if len(missingProviders) > maxMissingProviderCount {
						maxMissingProviderCount = len(missingProviders)
						maxProviderFilename = filename
					}
				}
				rootsToApply = append(rootsToApply, path.Dir(maxProviderFilename))
				rootProviders := lockfiles[maxProviderFilename].CompactProviders()
				logrus.Debugf("  applying %s", maxProviderFilename)
				logrus.Debugf("  because it has %d %s that are still missing. Its full set of providers: %v",
					maxMissingProviderCount,
					pluralize("provider", "providers", maxMissingProviderCount),
					rootProviders)
				for _, p := range rootProviders {
					if contains(appliedProviders, p) {
						continue
					}
					appliedProviders = append(appliedProviders, p)
				}
			}
			logrus.Infof("You can apply these providers with %d %s: %v",
				len(rootsToApply),
				pluralize("root", "roots", len(rootsToApply)),
				rootsToApply)
			for _, root := range rootsToApply {
				fmt.Println(root)
			}
			return nil
		},
	}

	return cmd
}

func newProviderHashesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hashes",
		Short: "Inspects all provider version hashes and notes exceptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			lockfiles, err := getLockfiles()
			if err != nil {
				return err
			}

			// map from provider name to version to hashes-hash to file list
			// ex: hashes["registry.terraform.io/hashicorp/aws"]["4.50.0"]["abcd...1234"] = ["terraform/roots/gold/500-regions/core/dev/us-west-1/lacework-integration/.terraform.lock.hcl"]
			hashes := make(map[string]map[string]map[string][]string)
			// lol maybe we need a custom data structure!

			for _, lf := range lockfiles {
				for _, p := range lf.Providers {
					if _, ok := hashes[p.ID]; !ok {
						hashes[p.ID] = make(map[string]map[string][]string)
					}
					if _, ok := hashes[p.ID][p.Version]; !ok {
						hashes[p.ID][p.Version] = make(map[string][]string)
					}
					hashOfHashes := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(p.Hashes, "\n"))))
					if _, ok := hashes[p.ID][p.Version][hashOfHashes]; !ok {
						hashes[p.ID][p.Version][hashOfHashes] = make([]string, 0)
					}
					hashes[p.ID][p.Version][hashOfHashes] = append(hashes[p.ID][p.Version][hashOfHashes], lf.Path)
				}
			}

			// print em out!
			// fmt.Printf("%+v\n", hashes)
			for providerID, versions := range hashes {
				fmt.Println(providerID)
				for version, hashes := range versions {
					fmt.Printf("\t%s\n", version)
					for hash, filenames := range hashes {
						fmt.Printf("\t\t%s: %d %s\n", hash, len(filenames), pluralize("file", "files", len(filenames)))
						if len(hashes) == 1 && !verbose && !vertrace {
							continue
						}
						for _, filename := range filenames {
							fmt.Printf("\t\t\t%s\n", filename)
						}
					}
				}
			}

			return nil
		},
	}

	return cmd
}

func newProviderVersionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "prints out all the versions required by lockfiles in or under the top directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			lockfiles, err := getLockfiles()
			if err != nil {
				return err
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
			logrus.Infof("Found %d provider %s", versionCount, pluralize("version", "versions", versionCount))
			fmt.Println(string(versionJSON))

			return nil
		},
	}

	return cmd
}

func newProviderWhyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "why PROVIDER[@VERSION]",
		Short: "prints out all the roots in or under the top directory that " +
			"require the given provider",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			lockfileNames, err := getLockfileNames()
			if err != nil {
				return err
			}

			target := args[0]
			targetParts := strings.Split(target, "@")
			targetProvider := targetParts[0]
			var targetVersion string
			if len(targetParts) > 1 {
				targetVersion = targetParts[1]
			}

			matches := make([]string, 0)
			for _, filename := range lockfileNames {
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

			logrus.Infof("%d %s found with the provider %s", len(matches), pluralize("root", "roots", len(matches)), target)

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

func getLockfileNames() ([]string, error) {
	ignore := append(ignoreNames, defaultIgnoreNames...)
	lockfileNames, err := findAll(".terraform.lock.hcl", topDir, ignore)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Found %d %s", len(lockfileNames), pluralize("lockfile", "lockfiles", len(lockfileNames)))
	return lockfileNames, nil
}

func getLockfiles() (map[string]*hcl.Lockfile, error) {
	lockfileNames, err := getLockfileNames()
	if err != nil {
		return nil, err
	}

	lockfiles := make(map[string]*hcl.Lockfile, 0)
	for _, filename := range lockfileNames {
		lf, err := hcl.ParseLockfile(filename)
		if err != nil {
			return nil, err
		}

		logrus.Debugf("%s (%d %s)", filename, len(lf.Providers), pluralize("provider", "providers", len(lf.Providers)))

		lockfiles[filename] = lf
	}
	return lockfiles, nil
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

func pluralize(single, plural string, count int) string {
	if count == 1 {
		return single
	}
	return plural
}
