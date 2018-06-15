# find-duplicates

**Build:** 
```
go build find-dups.go
```
**Usage:**
```
./find-dups --dir directory1 [--dir directory [--dir ...]]
```
---
#### This program is used to find duplicate files under all given directories.
**Notice:**
- Output of the program are the list of duplicate files grouped by --dir params.
- Duplicate files under the same --dir params would be shown the first one.
- Directories would be checked: 
    1. if valid (existing.)
    1. if duplicated, appointed multiple times by --dir,
        1. a/b/c/.. and a/b would be recognized as same, 
        1. links would be checked. 
