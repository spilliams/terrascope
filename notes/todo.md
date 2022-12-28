# TODO

## Up Next

- build a root module at all
- build a root module for only a subset of scopes
- build a root module with dependencies
- build a root module with a dependency outside the current scope
- build a root module in a scope with custom attributes
- cli should be able to add to an existing scope data file

unsorted shower thoughts:

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
3. `include(file, [scopeAddress], [attributesObj])` function for the root
   configuration
4. how should the CLI log?
   1. definitely to stdout, the way the user would like (e.g. verbose or quiet)
   2. definitely to a file in the build folder. How verbose?
   3. maybe also to a long-term cache in `~/.terraboots/logs/`. How verbose?
5. how much should cli care about scope value attributes?
   1. should `scope gen` ask for "what attributes does each scope value get?"
   2. should that info be embedded in the scope type definition in the project
      config?
6. for some commands, a `stat` subcommand might be nice. For `root`, it could show the number of scope matches it has.
7. for building, should we have a flag that states "please bootstrap". In other
   words, run all dependencies regardless of their last run time.
8. terraboots will want to inspect the last runtime of the roots. If a
   dependency has never been applied before, it should build it and apply it
   before any of its actual tasks.
9. using dependencies with custom attributes is a trap. The source of truth of
   "what are all my built roots" lies in the root config. If there's another
   root out there saying "actually add this flavor in too" it'll get hairy fast.
   So: when a root has a dependency on another root, it CANNOT allow for custom
   attribute values.
10. Track the source of an attribute value along with the value, for user
    debugging.

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

- I want to run arbitrary terraform commands:
  `terraboots tf my-root -- state mv 'module.a' 'module.b'`
- Terraboots should be able to generate a new monorepo project
- `terraboots root graph my-root` and/or `terraboots root graph`. With or
  without `--highlight-affected`.
