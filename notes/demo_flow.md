# demo flow

(checkboxes indicate the feature is built)

1. [x] create a new project from nothing.
   1. [x] prompts for top directory, project name, scope names, roots directory, scope data file name
   2. [x] generates a project-level terraboots.hcl
   3. [x] prompts for some sample scope values (this can also be its own command,
      something like `terraboots scope generate`)
   4. [x] generates a data.hcl
2. [x] create a new root from nothing.
   1. [x] prompts for root name, included scopes
   2. [x] generates a folder in the roots directory with `main.tf` and `terraboots.hcl`
3. [x] Building a root. Generally this should print out one or more folder names at
   the end.
   1. [x] `terraboots root build foo`
      builds a root named "foo", for all scopes
   2. [x] `terraboots root build foo "acme.gold.product.dev"`
      builds a root for a certain scope
   3. [x] `terraboots root build foo "acme.gold.*.*"`
      builds a root for scopes matching a given filter
4. [ ] Adding a new scope value to the existing data file

## Reactions

### Hello, world!

I'm running this by myself, on 14 Jan (commit dc96f3f7).

- [x] `project generate` "What scope types does your project use" could maybe use help text? At least pointing to a doc...nah, this will be on the Readme
- [x] `project generate` should say something in the end, like "terraboots.hcl file created"
- [x] `project generate-scopes` straight up bugged. It didn't like that the file didn't exist yet? Maybe we should ignore a "file not found" error in `readScopeData:97`
- [x] `scope show` panicked when i had an attribute in data.hcl that wasn't a string
- [ ] demo broke down after scope show and list, because I don't have a root generator yet

This was a great start I think! I have some bugs to iron out.
