# makeinvoice

## Requirements

### Go

* https://github.com/codegangsta/cli
* https://github.com/fatih/color

### External

* https://github.com/alanshaw/markdown-pdf

## TODO

- [ ] Handle CSS file.
- [ ] Support XLS document.
- [ ] Currently, we only import one CSV in one part of a document (`.Table`). Should be better to have multiple export points.
- [ ] Better template system (`{{ .Table }}` is not pretty).
- [ ] Vendor Go libraries.
- [ ] Remove the dependencies on [markdown-pdf](https://github.com/alanshaw/markdown-pdf).
