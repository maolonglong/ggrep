# g-grep

g-grep is a command-line tool that extends the functionality of the grep command by adding support for glob-style pattern matching.

## Installation

```bash
go install go.chensl.me/ggrep@latest
```

## Usage

```bash
$ echo 'foo\nbar\nbaz' | ggrep -Hn 'ba?'
(standard input):2:bar
(standard input):3:baz

$ ggrep -Hn 'ba?' ./hello.txt
./hello.txt:2:bar
./hello.txt:3:baz

$ ggrep -Hnv 'ba?' ./hello.txt
./hello.txt:1:foo

$ find . -type f -name '*.go' -exec ggrep -Hn '*redglob.*(*)' {} \;
./go.chensl.me/redglob/match_example_test.go:24:	fmt.Println(redglob.Match("foo", "f*"))
./go.chensl.me/ggrep/main.go:112:			matched := redglob.MatchBytes(l, pat)
```
