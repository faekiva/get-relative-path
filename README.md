# get-relative-path

Lets you find one path as relative to another. Figures out whether either the path written to or from is case sensitive (these can differ in cases like if a Linux NFS share is mounted to mac), and treat the paths as case sensitive if either are, or as case insensitive if neither are. Case sensitivity can also be specified manually.

Install with

```bash
go install github.com/faekiva/get-relative-path@latest
```


```bash
$ get-relative-path --help
Usage: get-relative-path [--relative-to RELATIVE-TO] [--case-sensitive CASE-SENSITIVE] [PATH]

Positional arguments:
  PATH                   if provided path is relative, it will be resolved relative to PWD first, then relative to the path provided with --relative-to

Options:
  --relative-to RELATIVE-TO [default: .]
  --case-sensitive CASE-SENSITIVE, -c CASE-SENSITIVE
                         options are true, false, or guess [default: guess]
  --help, -h             display this help and exit
```
