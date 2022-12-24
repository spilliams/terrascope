# demo flow

(checkboxes indicate the feature is built)

1. [ ] create a new project from nothing.
   1. [ ] prompts for top directory, project name, scope names, roots directory, scope data file name
   2. [ ] generates a project-level terraboots.hcl
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
4. [ ] Building a root. Generally this should print out one or more folder names at
   the end.
   1. [ ] `terraboots root build foo`
      builds a root named "foo", prompting along the way for what scope values
      to use (and offering a "press enter to leave blank and build for all
      subscopes).
   2. [ ] `terraboots root build foo --scopes "[acme.]gold.product.dev"`
      builds a root without prompting for scope values.
      The tree after running this build could be:

      ```txt
      └── terraform/roots/foo/
          ├── .terraboots/
          │   └── acme/gold/product/dev/
          │       ├── .terraboots.hcl
          │       ├── backend.tf
          │       ├── inputs.auto.tfvars
          │       ├── main.tf
          │       └── provider.tf
          ├── main.tf
          └── terraboots.hcl
      ```

      `.terraboots/` holds all the builds for the `foo` module.
      `.terraboots.hcl` contains all the scope data available to the module (for
      debugging)
      `backend.tf`, `provider.tf` and `inputs.auto.tfvars` will be generated
      from the module's `terraboots.hcl` configuration.

   3. [ ] `terraboots root build foo -a|--all`
      builds a root for all scope values, no prompting
5. [ ] Adding a new scope value to the existing data file
