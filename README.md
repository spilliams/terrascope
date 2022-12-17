# terraboots

Still very much a work in progress.

My attempt at a Terraform build orchestrator, for large platform projects with
hundreds of root modules.

This repository contains both the source code for the tool `terraboots`, as well
as a sample monorepo managed by that tool.

## Installation & Usage

Come back later :P

## Example Monorepo

The main configuration file for the monorepo is `terraboots.hcl`. This defines
where our root configurations live, as well as what `scope`s our monorepo deals
with.

Each of the roots has its own `terraboots.hcl` which contains configuration
details about that root. This includes what `scope`s apply to the root, and
what dependencies the root might have.

The main concept to terraboots is that it expects you to manage your platform
with many small terraform root configurations. You do this through `scope`s that
you define in the top level hcl, and selectively apply to each of your roots.

This allows you to maintain a few root "templates" that each could be planned
and applied dozens or hundreds of times depending on the permutations of your
scopes.

## Terraboots HCL

For schema documentation, please see `docs/hcl-schema.md`

## Terraboots CLI

See `src/README.md`.
