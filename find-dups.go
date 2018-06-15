package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Usage: ./main --dir dir1 --dir dir2 --dir ...

// FindDuplicates is a struct containing informations for finding process.
type FindDuplicates struct {
	Directories map[string]string   // used to regulate parameters of --dir, key: realpath value: path
	SizeFile    map[int64][]string  //used for recording all walked files by grouping them by file sizes
	Duplicates  map[string][]string //saves all duplicate files
}

type arrayFlags []string

// String is used for printing.
func (arr *arrayFlags) String() string {
	return fmt.Sprintf("%v", *arr)
}

// Set is used to regulate given paths to Directories
func (arr *arrayFlags) Set(value string) error {
	*arr = append(*arr, value)
	return nil
}

// InitDirs initialize the directories for scanning.
func (fd *FindDuplicates) InitDirs(arr arrayFlags) error {
	for _, n := range arr {
		_, err := os.Stat(n)
		if err != nil && os.IsNotExist(err) {
			return err
		}

		path, err := filepath.Abs(n)
		if err != nil {
			fmt.Printf("Error when get absolute path for %s: %s\n", n, err.Error())
			return err
		}

		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			fmt.Printf("Error when evaluate symbol link %s: %s\n", path, err.Error())
			return err
		}

		if dir, ok := fd.Directories[path]; ok {
			fmt.Printf("Warning: same directories: [%s, %s] => %s, use: %s\n", dir, n, path, n)
		}

		fd.Directories[path] = n
	}

	return nil
}

// ParseParams is used to parse parameters users input.
func ParseParams() []string {
	var arr arrayFlags
	flag.Var(&arr, "dir", "directory to compare")
	flag.Parse()

	return arr
}

// Traverse is used to walk through all files to group them by file sizes
func (fd *FindDuplicates) Traverse() {
	// CalculateSize calculate size of walked files
	CalculateSize := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			size := info.Size()
			if _, ok := fd.SizeFile[size]; !ok {
				fd.SizeFile[size] = []string{}
			}
			fd.SizeFile[size] = append(fd.SizeFile[size], path)
			return nil
		}

		return nil
	}

	for _, n := range fd.Directories {
		filepath.Walk(n, CalculateSize)
	}
}

// Summarize is used to group files of same size by md5sum.
func (fd *FindDuplicates) Summarize() {
	// Md5SumFile is ued to calculate file's md5sum
	Md5SumFile := func(file string) (string, error) {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		value := md5.Sum(data)
		return fmt.Sprintf("%x", value), nil
	}

	md5file := map[string][]string{}

	for _, files := range fd.SizeFile {
		if len(files) > 1 {
			for _, f := range files {
				r, _ := Md5SumFile(f)
				if _, ok := md5file[r]; !ok {
					md5file[r] = []string{}
				}
				md5file[r] = append(md5file[r], f)
			}
		}
	}

	for _, d := range fd.Directories {
		fd.Duplicates[d] = []string{}
	}

	for _, fs := range md5file {
		if len(fs) > 1 {
			mark := make(map[string]bool)
			for _, f := range fs {
				for _, d := range fd.Directories {
					if strings.HasPrefix(f, d) {
						if !mark[d] {
							fd.Duplicates[d] = append(fd.Duplicates[d], f)
							mark[d] = true
						}
						break
					}
				}
			}
		}
	}
}

// Report report duplicates with grouping them by directories
func (fd *FindDuplicates) Report() {
	for k, v := range fd.Duplicates {
		fmt.Printf("%s %v\n", k, v)
	}
}

func main() {
	dirs := ParseParams()
	fmt.Println(dirs)

	var findups = FindDuplicates{
		Directories: map[string]string{},
		SizeFile:    map[int64][]string{},
		Duplicates:  map[string][]string{},
	}

	findups.InitDirs(dirs)
	findups.Traverse()
	findups.Summarize()
	findups.Report()
}
