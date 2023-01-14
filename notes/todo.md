# TODO

## Up Next

- build a root module with dependencies
- build a root module with a dependency outside the current scope
- cli should be able to add to an existing scope data file
- `affected`

## Feature

- dependency detection (see `notes/root-dependencies.md`)
- People will probably want some kind of output wiring (not unlike Terragrunt's
  `dependency` blocks).
- Scope validation (e.g. "for the 'env' scope, values can only be 'dev',
  'stage', or 'prod'). See `notes/scope-validation.md`.

## CLI Commands

- Terraboots should be able to generate a new monorepo project
- `terraboots root graph my-root` and/or `terraboots root graph`. With or
  without `--highlight-affected`.
- `terraboots scope build SCOPE [-r|--recursive]` for building all roots for a
  single scope (and optionally all its children). Especially useful for
  gathering all the build directories of a whole domain
- `terraboots root clean ROOT` to delete temporary stuff
- `terraboots root show ROOT` to show root information

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
3. how should the CLI log?
   1. definitely to stdout, the way the user would like (e.g. verbose or quiet)
   2. definitely to a file in the build folder. How verbose? Also, how much to
      include? just the logs from buildContext or also get retroactive logs from
      main and project?
   3. maybe also to a long-term cache in `~/.terraboots/logs/`. How verbose?
4. how much should the cli care about scope value attributes?
   1. should `scope gen` ask for "what attributes does each scope value get?"
   2. should that info be embedded in the scope type definition in the project
      config?
5. for building, should we have a flag that states "please bootstrap". In other
   words, run all dependencies regardless of their last run time.
6. If a dependency has never been applied before, it should build it and apply
   it before any of its actual tasks.
7. Track the source of an attribute value along with the value, for user
   debugging.
8. project parser should make sure scopes have unique names with no special
   characters
9. how to get an exhaustive `affected` list?
   1. if a root module's source code changes (the `.tf` files), the whole root
      is affected. But maybe only the scopes that adopt the latest version?
   2. if a root module's config changes (`terraboots.hcl`), it depends
      1. if it's in a scope match, where some scopes are getting new attributes
         (which is something I haven't fully decided is a good idea yet), those
         scopes are affected.
   3. y'know what, for now lets keep it a little broad and say, if any files in
      the root (but not in `.terraboots/`) have changed, the whole root is
      affected for all matching scopes.
   4. scope data can change.
      1. If scope data changes, the whole scope could be affected. or maybe
         not, maybe only one sub-sub-sub scope uses a specific attribute.
      2. That, again, is oversolving. For now, we should hash each scope value,
         and save all the hashes in a file that's checked in. Then, when we're
         determining `affected` we compute a new hash, and figure out which
         scopes have changed.
      3. things to hash: all attributes in the scope
      4. also, detect if a scope is created or destroyed.
      5. if you don't want to check a hashfile into git, maybe the `affected`
         command clones the repo in a temp dir, runs a function to hash the
         clone, and runs the same function to hash the currect repo (or second
         clone).
      6. The other thing that can change (again, unimplemented) is wired outputs
         from root dependencies. Maybe what I need is a hash of the generated
         `.tfvars` file, and a matrix of scope and root.
         `map[scope][root]hash` turns into `map[scope][root]affected`?
         This would also solve the issue above, with parent scope attributes
         that are not used by all the child scopes.
      7. Yes, I'm more convinced now: the two things that can change in a root
         are the input values and the configuration itself. If the configuration
         changes (`.tf` files), the whole root is affected. If the input values
         change, certain scopes of the root are affected.
10. that brings up a new topic: what does it look like to destroy a scope?
    if we have to rename a domain, for instance, what does that entail?
    1. can we say "build all the roots for this scope and its children, and
       print out the list of build dirs" and then we can sequester those built
       directories for manual destruction? That seems plausible.
       1. that means we need a `terraboots scope build` command!
