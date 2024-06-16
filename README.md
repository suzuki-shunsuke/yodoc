# yodoc

Keep document including commands and their results up to date.

yodoc is a CLI to maintain documents including commands and their results.
When you write commands and their results in documents, it's hard to keep them up to date.
They would become outdated soon.

yodoc resolves this issue.
It executes commands, tests results, and generates documents based on template files.

## :warning: Security

yodoc executes commands according to the configuration file, so you shouldn't run yodoc with an untrusted configuration file.

## How to install

`yodoc` is a single binary written in Go.
So you only need to put the executable binary into `$PATH`.

```sh
go install github.com/suzuki-shunsuke/yodoc@latest
```

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

```

```

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
# File pathes are relative paths from the configuration file.
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

## Template

Template files are renderered by Go's [text/template](https://pkg.go.dev/text/template).

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

`dir` is renderered by Go's [text/template](https://pkg.go.dev/text/template).
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
  - expr: ExitCode != 0
```
</pre>

The annotation `#-yodoc check` must use at the top of a code block surrounded by <code>```</code>.
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

## LICENSE

[MIT](LICENSE)
