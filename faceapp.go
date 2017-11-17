package faceapp

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	ENDPOINT = "https://api.faceapp.io/api/v2.6"
)

type Filter string

const (
	FilterSmile     Filter = "smile"
	FilterSmile2    Filter = "smile_2"
	FilterHot       Filter = "hot"
	FilterOld       Filter = "old"
	FilterYoung     Filter = "young"
	FilterFemale    Filter = "female"
	FilterMale      Filter = "male"
	FilterBlack     Filter = "black"
	FilterCaucasian Filter = "caucasian"
	FilterAsian     Filter = "asian"
	FilterIndian    Filter = "indian"
	FilterFemale2   Filter = "female_2"
	FilterPan       Filter = "pan"
	FilterHitman    Filter = "hitman"
)

func (f Filter) NeedsCrop() bool {
	return f == FilterMale || f == FilterFemale ||
		f == FilterBlack || f == FilterCaucasian ||
		f == FilterAsian || f == FilterIndian ||
		f == FilterFemale2 || f == FilterPan ||
		f == FilterHitman
}

type Session struct {
	cli   http.Client
	appid string
}

func (s *Session) request(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "FaceApp/2.0.561 (Linux; Android 6.0)")
	req.Header.Set("X-FaceApp-DeviceID", s.appid)
	req.Header.Set("X-FaceApp-Priority", "1")
	return req, nil
}

func (s *Session) Init() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	s.cli.Jar = jar
	s.cli.Timeout = time.Second * 15
	var temp [8]byte
	if _, err := rand.Read(temp[:]); err != nil {
		return err
	}
	s.appid = hex.EncodeToString(temp[:])
	return nil
}

func (s *Session) UploadImage(in io.Reader) (string, error) {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	part, err := mw.CreateFormFile("file", "image.jpg")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, in)
	if err := mw.Close(); err != nil {
		return "", err
	}
	req, err := s.request("POST", ENDPOINT+"/photos", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	res, err := s.cli.Do(req)
	if err != nil {
		return "", err
	}
	var answer struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(res.Body).Decode(&answer); err != nil {
		res.Body.Close()
		return "", err
	}
	res.Body.Close()
	if answer.Code == "" {
		return "", errors.New("did not recognize image")
	}
	return answer.Code, nil
}

func (s *Session) GetImage(out io.Writer, code string, fil Filter, cropped bool) error {
	url := fmt.Sprintf(ENDPOINT+"/photos/%s/filters/%s?cropped=", code, fil)
	if cropped || fil.NeedsCrop() {
		url += "1"
	} else {
		url += "0"
	}
	res, err := s.request("GET", url, nil)
	if err != nil {
		return err
	}
	req, err := s.cli.Do(res)
	if err != nil {
		return err
	}
	io.Copy(out, req.Body)
	req.Body.Close()
	return nil
}
