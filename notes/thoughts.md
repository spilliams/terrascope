# thoughts

## The Core

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

## affected

what makes it so that a root is worth running?

1. the scope data changed
2. the root itself changed
3. the modules the root relied on changed? not sure how to do this quickly
   without either parsing the entire AST or assuming semver sources
