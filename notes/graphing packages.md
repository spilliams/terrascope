graphing packages

The current active uses of a graphing tool are directly to awalterschulze:

```text
internal/cli
    internal/hcl
        github.com/awalterschulze/gographviz
pkg/terrascope
    github.com/awalterschulze/gographviz
```

Then there's something off to the side a little bit, not imported by anyone (not
even a testing `main`):

```text
pkg/tfgraph
    github.com/spilliams/terrascope/internal/graphviz
        github.com/awalterschulze/gographviz
    github.com/spilliams/terrascope/pkg/grapher
```

If I had to guess (AND I DO), I was at the beginning of trying to write my own
replacement for awalterschulze...

ok, so what does it look like `grapher` is supposed to do?
oh gosh, and `pkg/interfaces.go`!
