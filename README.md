# distdin

Distribute stdin to multiple sub commands.

## Install

`go get github.com/fujiwara/distdin`

## Usage

```
$ distdin [-n Num of sub commands] [-v] command args ...
```

## Example: Distributed grep

```
$ distdin -n 4 grep Foo < file
```

`distdin` works as below.

1. Forks 4 grep commands.
2. Supplies STDIN for each line to each grep commands.

## LICENSE

The MIT License (MIT)

Copyright (c) 2016 FUJIWARA Shunichiro
