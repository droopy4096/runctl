# envctl

Manage multiple environment variable setups.

A lot of tools take their configuraetion from environment variables. For certain configurations several environment variables have to be set simulateneously and switching between configurations becomes cumbersome or one has to create wrapper scripts etc. This is where `envctl` comes to help.

## Define your environment variables configurations

create file `test.yaml`:

```yaml
myenv:
  - name: FOO
    value: foo 
  - name: BAR
    value: bar 
otherenv:
  - name: BAR
    value: barother
  - name: BAZ
    value: baz 
```

## Run any shell command within specified environment:

```shell
envctl -config-file test.yaml -command 'echo $BAR' -config myenv   
```
output is:

```
bar
```

```shell
envctl -config-file test.yaml -command 'echo $BAR' -config otherenv   
```
output is:

```
barother
```
