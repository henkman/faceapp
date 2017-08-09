package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

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
		string(faceapp.FilterBlack),
		string(faceapp.FilterCaucasian),
		string(faceapp.FilterAsian),
		string(faceapp.FilterIndian),
	}
)

var (
	_in      string
	_filters string
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&_in, "i", "", "jpg file input")
	flag.StringVar(&_filters, "f", "all", strings.Join(VALID_FILTERS, "|"))
	flag.Parse()
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}
	for v := range elements {
		encountered[elements[v]] = true
	}
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}

func main() {
	filterstrings := strings.Split(_filters, ",")
	filters := make([]string, 0, len(filterstrings))
	for _, s := range filterstrings {
		s = strings.TrimSpace(s)
		if s == "all" {
			filters = VALID_FILTERS
			break
		}
		for _, f := range VALID_FILTERS {
			if s == f {
				filters = append(filters, f)
			}
		}
	}
	filters = removeDuplicatesUnordered(filters)
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
	var wg sync.WaitGroup
	wg.Add(len(filters))
	for i, _ := range filters {
		go func(filter string) {
			fd, err := os.OpenFile(filepath.Join(dir, name+"_"+filter+".jpg"),
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := sess.GetImage(fd, code, faceapp.Filter(filter), false); err != nil {
				fd.Close()
				fmt.Println(err)
				return
			}
			fd.Close()
			wg.Done()
		}(filters[i])
	}
	wg.Wait()
}
