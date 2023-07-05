# Argo workflows expression tester


[Argo](https://github.com/argoproj/argo-workflows) has a complex [expression language](https://argoproj.github.io/argo-workflows/variables/#expression) to modify parameters in its workflows

It is difficult for a go novices to figure out how these work and what functions can be used.

This CLI provides a way to test the expressions before submitting the workflow to argo, it dumps failure information




## Installation

```bash
go install github.com/blacha/argo-expr@latest
```

## Usage

### Add 1 to a number

```bash
$ argo-expr "{{=asInt(input.parameters.name) + 1}}" --value input.parameters.name="1" 
2
```



### Create a sha256sum 

```bash
$ argo-expr '{{=sprig.sha256sum("hello world")}}' 
b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9

$ argo-expr '{{=sprig.sha256sum(input.value)}}' --value input.value="hello world" 
b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
```

### load from file

Read a JSON input file and compute the output

```json 
{
    "expression": "hello world 1+1:{{=1+1}} i:{{=i}}",
    "values": {
        "i": "4"
    }
}
```



```bash
$ argo-expr --from-file ./input.json
hello world 1+1:2 i:4
```

Both values and the expression from the file can be overridden with `--value` or expression.

```
$ argo-expr --from-file ./input.json --value i=1
hello world 1+1:2 i:2
```

Override input expression

```
$ argo-expr --from-file ./input.json "i:{{=asInt(i)+3}}"
i:7
```


### Error logs

When a template fails to run a somewhat helpful error message is displayed

```
$ argo-expr "{{=asInt('hello')}}" 

failed to evaluate expression: strconv.ParseInt: parsing "hello": invalid syntax (1:1)
 | asInt('hello')
 | ^
```