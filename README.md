# yodoc

Test command results and embed them into document.

> [!CAUTION]
> This project is still under development. Don't use this yet.

yodoc is a CLI to maintain documents.
When you write commands and their results in documents, it's hard to keep them fresh.
They would become stale.

`yodoc` enables you to run commands and embeds the result into documents.
Furthermore, you can test if the result is expected.

You can update tools automatically by Renovate, then you can also update documents automatically by `yodoc`.

## How to install

`yodoc` is a single binary written in Go. So you only need to put the executable binary into `$PATH`.

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
