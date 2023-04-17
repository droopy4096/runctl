# envctl

Manage multiple environment variable setups.

A lot of tools take their configuraetion from environment variables. For certain configurations several environment variables have to be set simulateneously and switching between configurations becomes cumbersome or one has to create wrapper scripts etc. This is where `envctl` comes to help.

## Define your environment variables configurations

create file `test.yaml`:

```yaml
myenv:
  - name: FOO
    value: foo 
    action: new
  - name: BAR
    value: bar 
    action: new
otherenv:
  - name: BAR
    value: barother
    action: replace
  - name: BAZ
    value: baz 
    action: merge
    separator: ":"
```

## Run any shell command within specified environment (`myenv`):

```shell
envctl -config-file test.yaml -config myenv -- 'echo $BAR'
```
output is:

```
bar
```

then run the same command using different environment configuration (`otherenv`)

```shell
envctl -config-file test.yaml -config otherenv -- 'echo $BAR'
```
output is:

```
barother
```

And of course one can use environment variables `ENVCTL_CONFIG` and `ENVCTL_ENV` to streamline use of `envctl` even further and have it as inobtrusive as possible:

```shell
export ENVCTL_CONFIG=test.yaml
export ENVCTL_ENV=myenv

envctl echo '$BAR'
```

### Action: new

only set variable if it hasn't been already set

### Action: replace

default behavior: replace existing variable value

### Action: merge

treat variable as array splitting using `separator`, append value to existing array

### Action: unset

unset existing value - as if it was never set to begin with
