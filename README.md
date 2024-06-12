# yodoc

Keep document including commands and their results up to date.

> [!CAUTION]
> This project is still under development. Don't use this yet.

yodoc is a CLI to maintain documents including commands and their results.
When you write commands and their results in documents, it's hard to keep them up to date.
They would become outdated soon.

yodoc resolves this issue.
It generates documents from a configuration file and templates.
It renderers templates by Go's [text/template](https://pkg.go.dev/text/template).
It supports custom template functions to execute commands while rendering documents.

e.g.

configuration file:

```yaml
tasks:
  - name: gh version
    action:
      run: gh version
```

template:

```
{{with Task "gh version"}}

{{.Command}}

{{.CombinedOutput}}

{{end}}
```

In the above example, yodoc runs `gh version` and embed the command and result into template and generate a document.
Furthermore, yodoc supports testing command results.

```yaml
tasks:
  - name: gh version
    action:
      run: gh version
    checks:
      - expr: Stdout startsWith "gh version"
```

## How to install

`yodoc` is a single binary written in Go.
So you only need to put the executable binary into `$PATH`.

```sh
go install github.com/suzuki-shunsuke/yodoc@latest
```

## How to use

```sh
yodoc run # Update document
yodoc watch # Watch file changes and update document automatically
```

You write template documents, then `yodoc` generates documents using them.
In template documents, you use Go's [text/template](https://pkg.go.dev/text/template).

```
{{Command "echo hello"}}
```

```
{{Task "yodoc version"}}
```

```
{{with Task "hello"}}	

$ {{.Command}}

{{.CombinedOutput}}

{{end}}
```

## Configuration file

yodoc.yaml

```sh
yodoc run src/README.md
```

```
<!-- yodoc: src/README.md -->
```

```yaml
# src/README.md => README.md
# src/src/README.md => src/README.md
src: src 
dest: .
files:
  - source: "*_yodoc.md"
  - dest: ""
tasks:
  - name: yodoc version
    command: yodoc -v
    checks:
      - exit_code: 0
```

## LICENSE

[MIT](LICENSE)
