# Fixtures

This directory is set up in a way similar to how I set up terraform repositories:

- `examples` is where I would keep example usages for my terraform modules. These examples often get used in the terraform tests.
- `modules` is where I keep individual modules
- `roots` is where my full terraform configurations live. This is what actually gets deployed
- `tests` is where I put code for unit, integration and end-to-end tests of examples and roots.
