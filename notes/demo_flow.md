# demo flow

(checkboxes indicate the feature is built)

1. [x] create a new project from nothing.
   1. [x] prompts for top directory, project name, scope names, roots directory, scope data file name
   2. [x] generates a project-level terraboots.hcl
   3. [x] prompts for some sample scope values (this can also be its own command,
      something like `terraboots scope generate`)
   4. [x] generates a data.hcl
2. [ ] create a new root from nothing.
   1. [ ] prompts for root name, included scopes
   2. [ ] generates a folder in the roots directory with `main.tf` and `terraboots.hcl`
3. [ ] create a root by importing one
   1. [ ] prompts for module address (e.g. GitHub URL), and included scopes
   2. [ ] infers module name from external source
   3. [ ] checks for name collisions
      1. [ ] resolves by offering to replace the current module
      2. [ ] offers to replace or keep existing terraboots.hcl (if applicable)
4. [x] Building a root. Generally this should print out one or more folder names at
   the end.
   1. [x] `terraboots root build foo`
      builds a root named "foo", for all scopes
   2. [x] `terraboots root build foo "acme.gold.product.dev"`
      builds a root for a certain scope
   3. [x] `terraboots root build foo "acme.gold.*.*"`
      builds a root for scopes matching a given filter
5. [ ] Adding a new scope value to the existing data file
