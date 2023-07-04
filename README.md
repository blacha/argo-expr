# Argo workflows expression tester


[Argo]() has a complex [expression language](https://argoproj.github.io/argo-workflows/variables/#expression) to modify parameters in its workflows

It is difficult for a go novice to figure out how these work.

This CLI provides a way to test the expressions before submitting the workflow to argo.


## Installation

```bash
go install https://github.com/blacha/argo-expr
```

## Usage

### simple math

```bash
argo-expr "{{=asInt(input.parameters.name) + 1}}" --value input.parameters.name="1"
```

output:
```json
{
  "input": "{{=asInt(input.parameters.name) + 1}}",
  "result": "2",
  "values": {
    "input.parameters.name": "1"
  }
}
```


### Create a sha256sum raw string

```bash
argo-expr '{{=sprig.sha256sum("hello world")}}'
```

output:
```json
{
  "input": "{{=sprig.sha256sum(\"hello world\")}}",
  "result": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
  "values": {}
}
```

### create a sha256sum hash from a input value

```bash
argo-expr '{{=sprig.sha256sum(input.value)}}' --value input.value="hello world"
```

output
```json
{
  "input": "{{=sprig.sha256sum(\"hello world\")}}",
  "result": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
  "values": {}
}