# yodoc

[![DeepWiki](https://img.shields.io/badge/DeepWiki-suzuki--shunsuke%2Fyodoc-blue.svg?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAyCAYAAAAnWDnqAAAAAXNSR0IArs4c6QAAA05JREFUaEPtmUtyEzEQhtWTQyQLHNak2AB7ZnyXZMEjXMGeK/AIi+QuHrMnbChYY7MIh8g01fJoopFb0uhhEqqcbWTp06/uv1saEDv4O3n3dV60RfP947Mm9/SQc0ICFQgzfc4CYZoTPAswgSJCCUJUnAAoRHOAUOcATwbmVLWdGoH//PB8mnKqScAhsD0kYP3j/Yt5LPQe2KvcXmGvRHcDnpxfL2zOYJ1mFwrryWTz0advv1Ut4CJgf5uhDuDj5eUcAUoahrdY/56ebRWeraTjMt/00Sh3UDtjgHtQNHwcRGOC98BJEAEymycmYcWwOprTgcB6VZ5JK5TAJ+fXGLBm3FDAmn6oPPjR4rKCAoJCal2eAiQp2x0vxTPB3ALO2CRkwmDy5WohzBDwSEFKRwPbknEggCPB/imwrycgxX2NzoMCHhPkDwqYMr9tRcP5qNrMZHkVnOjRMWwLCcr8ohBVb1OMjxLwGCvjTikrsBOiA6fNyCrm8V1rP93iVPpwaE+gO0SsWmPiXB+jikdf6SizrT5qKasx5j8ABbHpFTx+vFXp9EnYQmLx02h1QTTrl6eDqxLnGjporxl3NL3agEvXdT0WmEost648sQOYAeJS9Q7bfUVoMGnjo4AZdUMQku50McDcMWcBPvr0SzbTAFDfvJqwLzgxwATnCgnp4wDl6Aa+Ax283gghmj+vj7feE2KBBRMW3FzOpLOADl0Isb5587h/U4gGvkt5v60Z1VLG8BhYjbzRwyQZemwAd6cCR5/XFWLYZRIMpX39AR0tjaGGiGzLVyhse5C9RKC6ai42ppWPKiBagOvaYk8lO7DajerabOZP46Lby5wKjw1HCRx7p9sVMOWGzb/vA1hwiWc6jm3MvQDTogQkiqIhJV0nBQBTU+3okKCFDy9WwferkHjtxib7t3xIUQtHxnIwtx4mpg26/HfwVNVDb4oI9RHmx5WGelRVlrtiw43zboCLaxv46AZeB3IlTkwouebTr1y2NjSpHz68WNFjHvupy3q8TFn3Hos2IAk4Ju5dCo8B3wP7VPr/FGaKiG+T+v+TQqIrOqMTL1VdWV1DdmcbO8KXBz6esmYWYKPwDL5b5FA1a0hwapHiom0r/cKaoqr+27/XcrS5UwSMbQAAAABJRU5ErkJggg==)](https://deepwiki.com/suzuki-shunsuke/yodoc)

[Install](INSTALL.md)

Keep document including commands and their results up to date.

yodoc is a CLI to maintain documents including commands and their results.
When you write commands and their results in documents, it's hard to keep them up to date.
They would become outdated soon.

yodoc resolves this issue.
It executes commands, tests results, and generates documents based on template files.

## :warning: Security

yodoc executes commands based on template files, so you shouldn't run yodoc if they are untrusted.

## How to use

1. Scaffold a configuration file by `yodoc init`.

```sh
yodoc init # Scaffold a configuration file
```

2. Edit the configuration file.
3. Write templates in the source directory.
4. Generate documents by `yodoc run`.

```sh
yodoc run # Update document
```

By default, yodoc searches template files and process all template files, but you can also process only specific template files.

```sh
yodoc run README_yodoc.md # [...<file>]
```

## Usage

[USAGE.md](USAGE.md)

## Environment variables

- `YODOC_CONFIG`: a configuration file path
- `YODOC_LOG_COLOR`: log color mode: `auto`, `always`, `never`
- `YODOC_LOG_LEVEL`: trace, debug, info, warn, error, fatal, panic

## Configuration file paths

`yodoc` searches a configuration according to the following order.
If no configuration file is found, `yodoc` fails.

1. `--config`
1. `YODOC_CONFIG`
1. `yodoc.yaml`, `.yodoc.yaml` in the current directory :warning: `yodoc.yml` and `.yodoc.yml` are ignored

## Configuration file

yodoc.yaml

```yaml
# File paths are relative paths from the configuration file.
src: src # source directory where templates files are located
dest: . # destination directory where files are generated
delim:
  # delim is a pair of delimiters in templates.
  # This is optional.
  left: "{{" # The default value is "{{"
  right: "}}" # The default value is "}}"
```

`src` and `dest` are required.
yodoc searches template files from `src`, and generates documents in `dest`.
The file extension must be either `.md` or `.mdx`.

yodoc ignores directories `.git` and `node_modules`.

If `src` and `dest` are same, template file names must end with `_yodoc.md` or `_yodoc.mdx`, then yodoc generates `.md` or `.mdx`.

- `README_yodoc.md` => `README.md`
- `README_yodoc.mdx` => `README.mdx`

### JSON Schema

- [yodoc.json](json-schema/yodoc.json)
- https://raw.githubusercontent.com/suzuki-shunsuke/yodoc/refs/heads/main/json-schema/yodoc.json

If you look for a CLI tool to validate configuration with JSON Schema, [ajv-cli](https://ajv.js.org/packages/ajv-cli.html) is useful.

```sh
ajv --spec=draft2020 -s json-schema/yodoc.json -d yodoc.yaml
```

#### Input Complementation by YAML Language Server

[Please see the comment too.](https://github.com/szksh-lab/.github/issues/67#issuecomment-2564960491)

Version: `main`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/yodoc/main/json-schema/yodoc.json
```

Or pinning version:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/yodoc/v0.1.2/json-schema/yodoc.json
```

## Template

Template files are rendered by Go's [text/template](https://pkg.go.dev/text/template).

### Front matter

You can write [YAML Front matter](https://jekyllrb.com/docs/front-matter/) in the top of templates.

```md
---
dir: "{{.SourceDir}}"
dest: README.md
delim:
  left: "[["
  right: "]]"
---
```

Front matters are removed from generated document.
Front matter supports the following fields.

- `dir`
- `dest`
- `delim`

`dir` is rendered by Go's [text/template](https://pkg.go.dev/text/template).
The following variables are available.

- `SourceDir`: a directory where a template file exists
- `DestDir`: a directory where a file is generated
- `ConfigDir`: a directory where a configuration file exists

### Annotation

yodoc supports the following annotations.

- `#-yodoc hidden`: Executes a script and check if it succeeds but doesn't show the script and the output
- `#-yodoc run`: Executes a script and show the script. The command must succeed
- `#!yodoc run`: Executes a script and show the script. The command must fail
- `#-yodoc # `: Code comment. This isn't outputted
- `#-yodoc check`: Checks the result of the previous `#-yodoc run`
- `#-yodoc dir <dir>`: Change the directory where a command is executed

### `#-yodoc run`, `#!yodoc run`

Executes a script and show the script.
This annotation must use at the top of a code block surrounded by <code>```</code>.
You can use the result in templates.

e.g.

<pre>
```sh
#-yodoc run
npm test
```

```
{{.CombinedOutput -}}
```
</pre>

### `#-yodoc hidden`

Executes a script and check if it succeeds but doesn't show the script and the output.
This annotation must use at the top of a code block surrounded by <code>```</code>.
This annotation is used for preprocess, test, and clean up.

e.g.

<pre>
```sh
#-yodoc hidden
rm foo.json # Delete foo.json before `make init`
```

```sh
#-yodoc run
make init
```

```sh
#-yodoc hidden
test -f foo.json # Test if `make init` creates foo.json
```

# ...

```sh
#-yodoc hidden
rm foo.json # Clean up
```
</pre>

#### `#-yodoc dir <dir>`

Change the directory where a command is executed.
This annotation must use in a code block.

e.g.

```sh
#-yodoc hidden
#-yodoc dir foo
rm foo.json # Clean up
```

#### Checks

e.g.

<pre>
```sh
#-yodoc run
npm test
```

```yaml
#-yodoc check
checks:
  - expr: ExitCode == 0
```
</pre>

The annotation `#-yodoc check` is available at the top of a code block surrounded by <code>```</code>.
The content of the code block must be YAML.
`checks` is a list of checks.
All checks must be true, otherwise it fails to generate document.
check has the following fields.

- `expr`: An expression of [expr-lang/expr](https://github.com/expr-lang/expr).

The expression must return a boolean.
The following variables are passed.

- `Command`: Executed command
- `ExitCode`: Exit code
- `Stdout`: Standard output 
- `Stderr`: Standard error output
- `CombinedOutput`: Output combined Stdout and Stderr
- `RunError`: Error string if it fails to execute the command

The annotation `#-yodoc check` can be also used in `#-yodoc run` and `#!yodoc run` blocks.

```
#-yodoc check <expr>
```

<pre>
```sh
#-yodoc run
#-yodoc check Stdout contains "foo"
npm test
```
</pre>

### Template variables

In templates, the result of `#-yodoc run` is available.

- `Command`: Executed command
- `ExitCode`: Exit code
- `Stdout`: Standard output 
- `Stderr`: Standard error output
- `CombinedOutput`: Output combined Stdout and Stderr
- `RunError`: Error string if it fails to execute the command

e.g.

<pre>
```sh
#-yodoc run
npm test
```

```
{{.CombinedOutput}}
```
</pre>

### Template functions

[sprig functions](https://masterminds.github.io/sprig/) are available.
But the following functions are unavailable for security reason.

- `env`
- `expandenv`
- `getHostByName`

Furthermore, some custom functions are available.

- `Read`: Read a file

#### Read

```
{{Read "foo.yaml"}}
```

`Read` takes one argument, which is a file path to read.
The file path is a relative path from the template file.
`Read` returns a string which is the file content.

## Automatic update by CI

You can update documents by CI automatically, but you need to be care about the security.
We recommend executing yodoc by `push` event on the default branch and creating a pull request to update documents.

- https://github.com/suzuki-shunsuke/yodoc-workflow - Reusable workflow to update documents
- https://github.com/suzuki-shunsuke/poc-yodoc - Example of automation

## LICENSE

[MIT](LICENSE)
