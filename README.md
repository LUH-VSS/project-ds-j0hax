[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/IuYHBysz)


> [!NOTE]
> **Participant Information:**
> - Johannes Arnold, ErgebnisPIN `1022`

## Description

The project implements a very basic [MapReduce](https://static.googleusercontent.com/media/research.google.com/en//archive/mapreduce-osdi04.pdf) pattern.

The Mapper(s) read an arbitrary number of files in parallel, extract individual words, and send each unique word to a Reducer (determined by the word's hash modulo).

The Reducer(s) decode data from one or more Mappers, counting each unique word. On receiving a `SIGINT`, a Reducer will sort words and write them to a file with their respective counts.

### Installation and Usage

0. Compile the program with `go build -o <name>`
1. Compy the executable to the desired machines
2. Run `./<executable> map|reduce [options] [files]` to process the files. Help can be found by passing the `--help` flag.

## Changelog

### v0.2 (2024-01-26)
- Performance Improvements to address Prof. Dr. Rellermeyer's questions:
  - Use [encoding/gob](https://pkg.go.dev/encoding/gob) instead of JSON to reduce overhead
  - Stream over open TCP connection instead of multiple HTTP requests

### v0.1 (2024-01-25)
- Initial Release
