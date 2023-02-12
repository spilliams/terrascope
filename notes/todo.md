# TODO

## Up Next

- look for `TODO`s
- see `demo_flow` for bugs to sort out
- run the demo and figure out what's missing

- Track the source of an attribute value along with the value, for user debugging.
- `scope show` output looks weird. I feel like we should at least provide a
  `--json` option?
- root dependencies!
  - think real hard about it
  - build a root module with dependencies
  - build multiple root modules with the same dependency outside their current
    scope (e.g `acme.gold.commerce.dev`, `...stage` and `...prod` each depend on
    `acme.gold.networking.global`)
- graphing!
  - `scope graph` should build a graph of all compiled scopes
  - `scope graph FILTER` should build a graph of all filtered compiled scopes
    and their descendants. This could get weird for something like
    `acme.*.networking.*`, but cross that bridge later
  - `root graph` should build a graph of all the roots and their dependencies
- `terrascope root build --affected`
- pipeline for deploying `terrascope`. Version command, changelog, releases, etc.
- docs docs docs

## Feature

- dependency detection (see `notes/root-dependencies.md`)
- People will probably want some kind of output wiring (not unlike Terragrunt's
  `dependency` blocks).
- Scope validation (e.g. "for the 'env' scope, values can only be 'dev',
  'stage', or 'prod'). See `notes/scope-validation.md`.
- state migrations?

## CLI Commands

- `terrascope root graph my-root` and/or `terrascope root graph`. With or
  without `--highlight-affected`.
- `terrascope scope build SCOPE [-r|--recursive]` for building all roots for a
  single scope (and optionally all its children). Especially useful for
  gathering all the build directories of a whole domain
- `terrascope root clean ROOT` to delete temporary stuff
- `terrascope root show ROOT` to show root information
- `terrascope scope generate` to add to an existing scope data file

## unsorted shower thoughts

1. when adding a new scope to existing data files, do first impl by requiring a
   `--filename` flag.
   1. this way i can develop intelligent insertion later (e.g. "insert into the
      best file")
   2. if/when I introduce intelligent insertion, I can remove the requirement
      for the --filename flag. Backward compatibility!
2. along with `generate` for the backend and provider stuff, take another leaf
   out of Terragrunt's book: introduce a `terraform` block to use a module call.
   Then users can have their roots call on a versioned module.
   1. I wonder though, should this be a blue/green type of thing?
      Like, should I be able to say "for this `scopeMatch`, set these attributes
      to these values.
   2. Heck! combine this with Snap-like release channels, and you've got a
      really quick way to a stable-stable gold, stable-stable silver (prod
      replica), stable-candidate silver (prod release candidate), and
      candidate-stable silver (silver release candidate). Maybe I'm too far in
      the weeds here.
3. for building, should we have a flag that states "please run all dependendies".
4. If a dependency has never been applied before, it should build it and apply
   it before any of its actual tasks.
5. what does it look like to destroy a scope?
   if we have to rename a domain, for instance, what does that entail?
   1. can we say "build all the roots for this scope and its children, and
      print out the list of build dirs" and then we can sequester those built
      directories for manual destruction? That seems plausible.
      1. that means we need a `terrascope scope build` command!
