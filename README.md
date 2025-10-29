# Go-Editor

Another bare bone text editor written in Go.

## Installing termbox2
Install termbox2 in the editor folder

```bash
cd editor
git clone https://github.com/termbox/termbox2.git
cd termbox2
make
sudo make install
```

## Build

```bash
go build -o go_editor main.go
```

## Run

```bash
./go_editor
```