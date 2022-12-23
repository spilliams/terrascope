# TODO

## Up Next

- read a scope data file
- cli should be able to add to an existing scope data file

unsorted shower thoughts:

1. scope data attributes should be by themselves, not in an object
2. i want to support having multiple scope data files
   1. when adding a new scope to existing data files, do first impl by requiring
      a --filename flag.
      1. this way i can develop intelligent insertion later (e.g. "insert into
         the best file")
      2. if/when I introduce intelligent insertion, I can just remove the
         requirement for the --filename flag. Backwards compatible!
3. scope value addresses could follow terraform's example, with the types
   embedded. e.g. `org.acme.platform.gold.domain.core`
   1. do I want to take this a step more generic and use
      `module.acme.module.gold.module.core`? Unsure
4. where is the best place to put the "this scope address gets these roots"
   data?
   1. could keep it in scope data file where it is now: it *feels* like an
      extension of the scopes.
   2. could put it in the root configs!
      1. This could be especially powerful if we use regex to fill it out:
         The `root` block could have any number of `scopeMatch` sub blocks,
         that use an attribute named `addressPattern`. Example pattern:
         `org.acme.platform.*.domain.*.environment.[dev|stage|prod]`
      2. removing `roots` from the scope data means that we no longer have any
         attributes in the scope blocks that are not user-defined! KISS
      3. using a process that allows for regex means we can reduce the
         copy-paste there too!
5. along with `generate` for the backend and provider stuff, take another leaf
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
6. `include(file, [scopeAddress], [attributesObj])` function for the root
   configuration
7. how should the CLI log?
   1. definitely to stdout, the way the user would like (e.g. verbose or quiet)
   2. definitely to a file in the build folder. How verbose?
   3. maybe also to a long-term cache in `~/.terraboots/logs/`. How verbose?

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
