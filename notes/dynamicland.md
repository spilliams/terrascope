(split from Obsidian note cards/projects/Dynamicland-RealtalkOS.md)

Imagine a workspace like Realtalk, for building cloud infrastructure.

## Physical computer

The physical computer of it would have these components:

- work surface(s) to put objects
- projector(s) to display computation on the work surfaces (output)
- camera(s) to watch the work surfaces (input)
- objects themselves represent infrastructure resources (e.g. `provider`, `backend`, `resource`, `data`), processes (`plan`, `apply`), or even whole systems (`terraform`, `cloudformation`, `aws cli`).
- scanner(s) to read and understand objects as the humans have made them. Maybe these are "really good cameras", maybe actual flatbed scanners, maybe "smartphone take a photo and upload it".
- printer to have computer draw on the objects (ID codes)

## Use cases

Imagine a scenario of using this:

1. print out a token
2. configure the token's type and identifier (hcl block type and labels)
3. place the token on the board. Computer responds by asking for inputs and showing you outputs.
4. You can use a tool to inspect any of the inputs and outputs, and then connect them to things.
5. You can use a tool to run `terraform` commands on the configuration, and you'll see output in the form of log text and tokens highlighted with status color.

Once you have the configuration fairly stable, you could print out 4 corner-tokens, and place them around your objects to form a module, then draw out the variables and outputs of the module. At this point you could use a tool to "publish" the module to a github tag or to a terraform registry.

With many objects bracketed by the 4 corner-tokens of a module definition, you could sweep them off the board and replace them with a smaller "module" object to represent them all. If you want to edit the inside of the module? ask the computer to project the module contents back, and place your module's internal objects to match, then proceed with the edits. Save and publish a new version.

Using the "Inspector" tool on a module block could show details like readme, version constraints, changelog, etc. And the "Wire" tool could connect inputs, outputs, providers, etc.

## Virtual dynamic space

This system could probably have a web interface too, so collaborators can have their own workspaces digitally. The digital workspace could transclude a camera view of any physical workspace and vice-versa.
