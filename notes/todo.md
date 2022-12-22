# TODO

## Up Next

- cli should be able to read a project config and create a skeletal scope data file
- read a scope data file

## Feature

- `affected`
- Templating
- dependency detection (see `notes/root-dependencies.md`)
- People will probably want some kind of output wiring (not unlike Terragrunt's
  `dependency` blocks).
- Scope validation (e.g. "for the 'env' scope, values can only be 'dev',
  'stage', or 'prod'). See `notes/scope-validation.md`.

## Bulletproofing

- project parser should make sure scopes have unique names with no special
  characters

## CLI Commands

- I want to run arbitrary terraform commands: `terraboots tf my-root -- state mv 'module.a' 'module.b'`
- Terraboots should be able to generate a new monorepo project
- `terraboots root graph my-root` and/or `terraboots root graph`. With or
  without `--highlight-affected`.
