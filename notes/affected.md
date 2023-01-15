# "Affected"

What makes it so that a root is worth running?

1. the scope data changed
2. the root itself changed
3. the modules the root relied on changed? not sure how to do this quickly
   without either parsing the entire AST or assuming semver sources

How to get an exhaustive `affected` list?

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
