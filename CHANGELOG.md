# Changelog

See [Versioning](./README.md#Versioning) for how to regard these version numbers.

## 0.2.0

- The command `terrascope root graph-dependencies` is now
  `terrascope project graph-roots`, to better align with other command naming
  patterns.

## 0.1.2

- Adds a command `terrascope root show ROOT` which will print out the terrascope
  configuration file for the given root.
- Generated context files will now show which scope each value came from, for
  debugging purposes.
- Bumps golang.org/x/text from 0.3.7 to 0.3.8

## 0.1.1

- Adds a new command: `terrascope root clean ROOT [SCOPE]`, to clean up builds.
-nAlso adds a shorthand to the `--dry-run` flag: `-d`.

## 0.1.0

Initial release of the CLI tool and package!

The features in this release are mainly:

- generate a new project
- add scope data to the project
- generate a new root
- build a root for a set of scopes

There is unstable support for running terraform commands in a root's scopes.
