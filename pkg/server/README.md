# server

This package provides a webserver.

The router uses the [gin framework](github.com/gin-gonic/gin).

The graphing and visualization on the frontend are done with [graphology](https://graphology.github.io/) and [sigmajs](https://www.sigmajs.org/), repsectively.

The initial goal is for the backend to manage a real terraform configuration. The frontend will show a model of the configuration to the user, provide the user ways to manipulate the configuration (add, remove, and edit blocks, run `terraform` commands), and return logs and warnings from `terraform`.

The second major milestone is getting this interface and model to represent a terrascope configuration. That is, the user should be able to modify the scope types and values, as well as visualize the state of all the scope-configurations.
