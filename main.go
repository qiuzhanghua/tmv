package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

/*
Usage
mkdir -p  src/main/java/cn/com/taiji/sample
touch src/main/java/cn/com/taiji/sample/App.java
mkdir -p  src/test/java/cn/com/taiji/sample
touch src/test/java/cn/com/taiji/sample/Test.java
rm -rf src/main/java/com
rm -rf src/test/java/com

go run *.go src/main/java/cn/com/taiji/sample src/main/java/com/example

*/

var regex = regexp.MustCompile("[$][{].+?[}]|[%].+?[%]")

func main() {
	if len(os.Args) <= 2 {
		fmt.Printf("Usage:\n   %s <from> <to>\n", os.Args[0])
		os.Exit(0)
	}
	fromOrig := ReplaceEnvString(os.Args[1])
	toOrig := ReplaceEnvString(os.Args[2])
	dir, _ := os.Getwd()
	if len(os.Args) == 4 {
		dir = os.Args[3]
	}
	if fromOrig == toOrig {
		os.Exit(0)
	}
	fromList := strings.Split(fromOrig, "/\\")
	from := Join(fromList)
	toList := strings.Split(toOrig, "/\\")
	to := Join(toList)
	mv(from, to, dir)
	os.Exit(0)
}

//mv
func mv(from string, to string, base string) {
	sep := "/"
	if runtime.GOOS == "windows" {
		sep = "\\"
	}
	baseAbs, err := AbsPath(base)
	AssetNil(err)
	fromList := strings.Split(from, sep)
	toList := strings.Split(to, sep)
	src := JoinVar(base, from)
	dest := JoinVar(base, to)
	fi, err := os.Stat(src)
	AssetNil(err)
	if fi.IsDir() {
		tempDir, err := ioutil.TempDir(base, "tdp_")
		AssetNil(err)
		err = os.Remove(tempDir)
		AssetNil(err)
		err = os.Rename(src, tempDir)
		AssetNil(err)
		destParent := Join(toList[:len(toList)-1])
		err = os.MkdirAll(JoinVar(base, destParent), os.ModePerm)
		AssetNil(err)
		err = os.Rename(tempDir, dest)
		AssetNil(err)
	} else {
		tempFile, err := ioutil.TempFile(base, "tdp_")
		AssetNil(err)
		filename := tempFile.Name()
		err = os.Remove(filename)
		AssetNil(err)
		err = os.Rename(src, filename)
		AssetNil(err)
		destParent := Join(toList[:len(toList)-1])
		err = os.MkdirAll(JoinVar(base, destParent), os.ModePerm)
		AssetNil(err)
		err = os.Rename(filename, dest)
		AssetNil(err)
	}
	for i := 1; i < len(fromList); i++ {
		p := Join(fromList[0 : len(fromList)-i])
		p = JoinVar(base, p)
		p, err = AbsPath(p)
		AssetNil(err)
		if p == baseAbs { // 到了不应删除的位置
			break
		}
		isEmpty, err := IsEmpty(p)
		AssetNil(err)
		if !isEmpty {
			break
		}
		err = os.Remove(p)
		AssetNil(err)
	}
}

func Join(path []string) string {
	sep := "/"
	if runtime.GOOS == "windows" {
		sep = "\\"
	}
	return strings.Join(path, sep)
}

func JoinVar(path ...string) string {
	sep := "/"
	if runtime.GOOS == "windows" {
		return strings.Join(path, "\\")
	}
	return strings.Join(path, sep)
}

func ReplaceEnvString(s string) string {
	envs := regex.FindAllString(s, -1)
	if len(envs) == 0 {
		return s
	}
	for _, e := range envs {
		if len(e) < 3 {
			continue
		}
		var e2 string
		if strings.HasPrefix(e, "$") {
			e2 = e[2 : len(e)-1]
		} else {
			e2 = e[1 : len(e)-1]
		}
		env := os.Getenv(e2)
		s = strings.ReplaceAll(s, e, env)
	}
	return s
}

func AbsPath(p string) (string, error) {
	absPath, err := filepath.EvalSymlinks(p)
	if err != nil {
		return ".", err
	}
	absPath, err = filepath.Abs(absPath)
	if err != nil {
		return ".", err
	}
	return absPath, nil
}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func AssetNil(err interface{}) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
