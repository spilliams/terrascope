# terraboots

Still very much a work in progress.

My attempt at a Terraform build orchestrator, for large infrastructure projects
with hundreds of root modules.

This repository contains both the source code for the tool `terraboots`, as well
as a sample monorepo managed by that tool.

## The Big Idea

Here's the core of terraboots:

- you don't need anything other than terraform to run it. This means its config
  files are dead simple.
- it will tell you when it is improperly configured. This is important because
  configuration is confusing! And guardrails are just as important as docs.
  This means the configs will all be in HCL, and the tool will validate them at
  runtime.
- it can solve the `affected` problem quickly. This means it has to be
  performant. And maybe it has to be extensible, because people have different
  ideas of what it means to be "affected".
- it can solve root dependencies quickly to generate a parallelizable task
  manifest. Again, it has to be performant.

The primary goal of `terraboots root list` is to make sure we can build a matrix
for GitHub Actions or any other CI runner to know exactly what to run. It might
want to take root dependency into account too.

## Installation & Usage

Come back later :P

## Example Monorepo

The main configuration file for the monorepo is `terraboots.hcl`. This defines
where our root configurations live, as well as what `scope`s our monorepo deals
with.

Each of the roots has its own `terraboots.hcl` which contains configuration
details about that root. This includes what `scope`s apply to the root, and
what dependencies the root might have.

The main concept to terraboots is that it expects you to manage your
infrastructure with many small terraform root configurations. You do this through `scope`s that
you define in the top level hcl, and selectively apply to each of your roots.

This allows you to maintain a few root "templates" that each could be planned
and applied dozens or hundreds of times depending on the permutations of your
scopes.

## Terraboots HCL

For schema documentation, please see `docs/hcl-schema.md`

## Terraboots CLI

See `src/README.md`.
