# UmlGen

A tool for generating UML Class Diagrams from JSON textual representations.

Created during Knight Hacks 2021.

## Installation

```shell
git clone git@github.com:cass-dlcm/umlgen.git
cd umlgen
go build umlgen.go
```

## Use

A few examples of use are shown here.

Using piping and output redirection:
```shell
cat input.json | umlgen >> diagram.svg
```

Using the `-input` and `-output` flags:
```shell
umlgen -input input.json -output diagram.svg
```