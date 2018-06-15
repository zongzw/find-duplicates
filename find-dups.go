package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ./main --dir dir1 --dir dir2

// Directories is type of map used for regulating --dir parameters, key: realpath value: path
type Directories map[string]string

// SizeFile is a map used for recording all walked files by grouping them by file sizes
var SizeFile = make(map[int64][]string)

// DupFiles saves all duplicate files
var DupFiles = [][]string{}

// String is used for printing.
func (d *Directories) String() string {
	rlt := []string{}
	for _, v := range *d {
		rlt = append(rlt, v)
	}

	return fmt.Sprint(rlt)
}

// Set is used to regulate given paths to Directories
func (d *Directories) Set(value string) error {
	info, err := os.Stat(value)
	if err != nil && os.IsNotExist(err) {
		return err
	}

	path, err := filepath.Abs(value)
	if err != nil {
		fmt.Printf("Error when get absolute path for %s: %s\n", value, err.Error())
		return err
	}
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		path, err := filepath.EvalSymlinks(path)
		if err != nil {
			fmt.Printf("Error when evaluate symbol link %s: %s\n", path, err.Error())
			return err
		}
	}

	if dir, ok := (*d)[path]; ok {
		fmt.Printf("Given path %s is same as %s, which is actually %s, overriden.\n", value, dir, path)
	}

	(*d)[path] = value
	return nil
}

// CalculateSize calculate size of walked files
func CalculateSize(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		size := info.Size()
		if _, ok := SizeFile[size]; !ok {
			SizeFile[size] = []string{}
		}
		SizeFile[size] = append(SizeFile[size], path)
		return nil
	}

	return nil
}

// Md5SumFile is ued to calculate file's md5sum
func Md5SumFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	value := md5.Sum(data)
	return fmt.Sprintf("%x", value), nil
}

func main() {
	var ds Directories
	flag.Var(&ds, "dir", "directory to compare")

	flag.Parse()
	fmt.Printf("%v\n", ds)

	for _, n := range ds {
		filepath.Walk(n, CalculateSize)
	}

	/*
		for _, fs := range SizeFile {
			if len(fs) > 1 {
				m := map[string][]string{}
				for _, f := range fs {
					r, _ := Md5SumFile(f)
					if _, ok := m[r]; !ok {
						m[r] = []string{}
					}
					m[r] = append(m[r], f)
				}
				for _, v := range m {
					if len(v) > 1 {
						DupFiles = append(DupFiles, v)
					}
				}
			}
		}

		dd := []string{}
		for _, n := range DupFiles {
			lastDir := ds[len(ds)-1]
			for _, f := range n {
				if strings.HasPrefix(f, lastDir) {
					dd = append(dd, f)
				}
			}
		}

		fmt.Println("")
		for _, n := range dd {
			fmt.Printf("%s\n", n)
		}
	*/
}
