package faceapp

import (
	"os"
	"testing"
)

func TestSimple(t *testing.T) {
	var s Session
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	var code string
	{
		fd, err := os.Open("test/bill.jpg")
		if err != nil {
			t.Fatal(err)
		}
		temp, err := s.UploadImage(fd)
		fd.Close()
		if err != nil {
			t.Fatal(err)
		}
		code = temp
	}
	{
		for _, fil := range []Filter{
			FilterSmile,
			FilterSmile2,
			FilterHot,
			FilterOld,
			FilterYoung,
			FilterFemale,
			FilterMale,
			FilterBlack,
			FilterCaucasian,
			FilterAsian,
			FilterIndian,
		} {
			fd, err := os.OpenFile("test/bill_"+string(fil)+".jpg",
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0750)
			if err != nil {
				t.Fatal(err)
			}
			if err := s.GetImage(fd, code, fil, false); err != nil {
				fd.Close()
				t.Fatal(err)
			}
			fd.Close()
		}
	}
}
