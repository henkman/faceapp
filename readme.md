faceapp
=======

small package for faceapp

```go
import (
	"os"
	"github.com/henkman/faceapp"
)

func main() {
	var s faceapp.Session
	if err := s.Init(); err != nil {
		panic(err)
	}
	var code string
	{
		fd, err := os.Open("test/bill.jpg")
		if err != nil {
			panic(err)
		}
		temp, err := s.UploadImage(fd)
		fd.Close()
		if err != nil {
			panic(err)
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
				panic(err)
			}
			if err := s.GetImage(fd, code, fil, false); err != nil {
				fd.Close()
				panic(err)
			}
			fd.Close()
		}
	}
}
```

original: ![alt original](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill.jpg "original")

smile: ![alt smile](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_smile.jpg "smile")

smile_2: ![alt smile_2](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_smile_2.jpg "smile_2")

hot: ![alt hot](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_hot.jpg "hot")

old: ![alt old](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_old.jpg "old")

female: ![alt female](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_female.jpg "female")

male: ![alt male](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_male.jpg "male")

black: ![alt black](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_black.jpg "black")

caucasian: ![alt caucasian](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_caucasian.jpg "caucasian")

asian: ![alt asian](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_asian.jpg "asian")

indian: ![alt indian](https://raw.githubusercontent.com/henkman/faceapp/master/test/bill_indian.jpg "indian")
