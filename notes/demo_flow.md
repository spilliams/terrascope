# demo flow

(checkboxes indicate the feature is built in feature/scopes)

1. [x] create a new project from nothing.
   1. [x] prompts for top directory, project name, scope names, roots directory,
      scope data file name
   2. [x] generates a project-level terrascope.hcl
   3. [x] prompts for some sample scope values (this can also be its own command,
      something like `terrascope scope generate`)
   4. [x] generates a data.hcl
2. [x] create a new root from nothing.
   1. [x] prompts for root name, included scopes
   2. [x] generates a folder in the roots directory with `main.tf` and
      `terrascope.hcl`
3. [x] Building a root. Generally this should print out one or more folder names
   at the end.
   1. [x] `terrascope root build foo`
      builds a root named "foo", for all scopes
   2. [x] `terrascope root build foo "acme.gold.product.dev"`
      builds a root for a certain scope
   3. [x] `terrascope root build foo "acme.gold.*.*"`
      builds a root for scopes matching a given filter
4. [ ] Adding a new scope value to the existing data file
