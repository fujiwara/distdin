# distdin

Distribute stdin to multiple sub commands.

## Usage

```
$ distdin [-n Num of sub commands] [-v] command args ...
```

## Example: Distribulted grep

```
$ distdin -n 4 grep Foo < file
```

`distdin` works as below.

1. Forks 4 grep commands.
2. Supplies STDIN for each line to each grep commands.

## LICENSE

The MIT License (MIT)

Copyright (c) 2016 FUJIWARA Shunichiro
