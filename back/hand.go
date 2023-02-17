package main

import (
	"os"

	"github.com/qwertyqq2/filebc/files"
)

func getData() ([][]byte, error) {
	data1, err := os.ReadFile("htmlfiles/htmlExample1.html")
	if err != nil {
		return nil, err
	}
	data2, err := os.ReadFile("htmlfiles/htmlExample2.html")
	if err != nil {
		return nil, err
	}
	data3, err := os.ReadFile("htmlfiles/htmlExample3.html")
	if err != nil {
		return nil, err
	}
	return [][]byte{data1, data2, data3}, nil
}

func newColl() (*files.Collector, error) {
	coll, err := files.NewCollector("uname")
	if err != nil {
		return nil, err
	}
	data, err := getData()
	if err != nil {
		return nil, err
	}
	for _, d := range data {
		f := files.NewFile(string(d))
		if err := coll.InsertFile(f); err != nil {
			return nil, err
		}
	}
	return coll, nil
}

func Get(coll *files.Collector) ([]string, error) {
	fsd, err := coll.LDB().GetFiles()
	if err != nil {
		return nil, err
	}
	posts := make([]string, len(fsd))
	for _, fs := range fsd {
		fssrt, err := files.Deserialize(fs)
		if err != nil {
			return nil, err
		}
		posts = append(posts, string(fssrt.Data))
	}
	return posts, nil
}
