# RANDOMART

## Build

```sh
go build -o .out/randomart
```

Executable binary at: `.out/randomart`

## Usage

```sh
randormart -a sha256 -w 29 -h 15 -c true -k ~/.ssh/id_ed25519.pub
```

```sh
randormart -a sha256 -w 29 -h 15 -c true "hello"
```

Output:

```asciiarmor
+[SHA256] [2cf24dba5fb0a30e]+
+-----------------+
|   +.+..         |
|E...o +          |
|+..o.+ o         |
|+.+.o *.         |
| = +.o. S        |
|..* .o = o       |
|.=..o o + .      |
|=.+o . o o       |
|*O+  .+..        |
+-----------------+
```

