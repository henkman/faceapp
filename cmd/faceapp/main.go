package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/henkman/faceapp"
)

var (
	VALID_FILTERS = []string{
		string(faceapp.FilterSmile),
		string(faceapp.FilterSmile2),
		string(faceapp.FilterHot),
		string(faceapp.FilterOld),
		string(faceapp.FilterYoung),
		string(faceapp.FilterFemale),
		string(faceapp.FilterMale),
		"all",
	}
)

var (
	_in      string
	_filters string
)

func init() {
	flag.StringVar(&_in, "i", "", "jpg file input")
	flag.StringVar(&_filters, "f", "all", strings.Join(VALID_FILTERS, "|"))
	flag.Parse()
}

func main() {
	filterstrings := strings.Split(_filters, ",")
	filters := make([]faceapp.Filter, 0, len(filterstrings))
	for _, s := range filterstrings {
		s = strings.TrimSpace(s)
		for _, f := range VALID_FILTERS {
			if s == f {
				filters = append(filters, faceapp.Filter(s))
			}
		}
	}
	if _in == "" || len(filters) == 0 {
		flag.Usage()
		return
	}
	var sess faceapp.Session
	if err := sess.Init(); err != nil {
		fmt.Println(err)
		return
	}
	fd, err := os.Open(_in)
	if err != nil {
		fmt.Println(err)
		return
	}
	code, err := sess.UploadImage(fd)
	fd.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	dir := filepath.Dir(_in)
	file := filepath.Base(_in)
	ext := filepath.Ext(file)
	name := file[0 : len(file)-len(ext)]
	for _, filter := range filters {
		fd, err := os.OpenFile(filepath.Join(dir, name+"_"+string(filter)+".jpg"),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := sess.GetImage(fd, code, filter, false); err != nil {
			fd.Close()
			fmt.Println(err)
			return
		}
		fd.Close()
	}
}
