# yodoc

Keep document including commands and their results up to date.

> [!CAUTION]
> This project is still under development. Don't use this yet.

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
tasks:
  # A list of tasks
  # ...
  - name: gh version
    action:
      run: gh version
```

`src`, `dest`, and `tasks` are required.
yodoc searches template files from `src`, and generates documents in `dest`.
The file extension must be either `.md` or `.mdx`.

### .tasks

```yaml
- name: gh version
  action:
    run: gh version
  # before and after are optional.
  # The structure is same as `action`.
  before:
    run: mkdir foo
  after:
    run: rm -R foo
  checks:
    # A list of checks.
    - expr: ExitCode == 1
```

The command of `before` and `after` must succeed, otherwise it fails to render the template.

### .tasks[].action, before, after

```yaml
# Either `run` or `script` is required.
# Others are optional.
run: gh version
script: foo.sh # relative path from the configuration file
shell: ["sh", "-c"] # shell to execute the command
dir: foo # The directory where the command is executed. A relative path from the configuration file
env:
  # environment variables to pass the command
  # Go's template is available.
  FOO: '{{env "HOME"}}/foo'
```

### .tasks[].checks[]

```yaml
# expr is an expression of 
- expr: ExitCode == 0
```

`expr` is an expression of [expr-lang/expr](https://github.com/expr-lang/expr).
The expression must return a boolean.

The following variables are passed.

- `Command`: Executed command
- `ExitCode`: Exit code
- `Stdout`: Standard output 
- `Stderr`: Standard error output
- `CombinedOutput`: Output combined Stdout and Stderr
- `RunError`: Error string if it fails to execute the command

All checks must be true, otherwise the task fails.

## Template

Template files are renderered by Go's [text/template](https://pkg.go.dev/text/template).

### Template functions

[sprig functions](https://masterminds.github.io/sprig/) are available.
But the following functions are unavailable for security reason.

- `env`
- `expandenv`
- `getHostByName`

Furthermore, some custom functions are available.

- `Read`: Read a file
- `Task`: Execute a task

#### Read

```
{{Read "foo.yaml"}}
```

`Read` takes one argument, which is a file path to read.
The file path is a relative path from the template file.
`Read` returns a string which is the file content.

#### Task

```
{{with Task "gh version"}}

{{.Command}}

{{.CombinedOutput}}

{{end}}
```

Task takes one argument, which is a task name.
The task must be defined in the configuration file.

```yaml
tasks:
  - name: gh version
    # ...
```

1. Run `before`
1. Run `action`
1. Run `after`
1. Test the result by `checks`

Task executes a command and returns its result.
The result has the following attributes.

- `Command`: Executed command
- `ExitCode`: Exit code
- `Stdout`: Standard output 
- `Stderr`: Standard error output
- `CombinedOutput`: Output combined Stdout and Stderr
- `RunError`: Error string if it fails to execute the command

## Automatic update by CI

You can update documents by CI automatically, but you need to be care about the security.
We recommend executing yodoc by `push` event on the default branch and creating a pull request to update documents.

## LICENSE

[MIT](LICENSE)
