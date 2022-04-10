package service

import (
	"bytes"
	"cctv-image-receiver/conf"
	"cctv-image-receiver/util"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

type ImageJob struct {
	CctvCountMap map[string]int
	ImageOutPath string
	Cctv         conf.Cctv
}

func (this ImageJob) Run() {
	getImage(this.Cctv, this.ImageOutPath, this.CctvCountMap)
}

//取得影像
func getImage(cctv conf.Cctv, imageOutPath string, cctvCountMap map[string]int) {

	url := util.CombineString("http://", cctv.Ip, cctv.ImageUrlSuffix)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", "Basic "+util.BasicAuth(cctv.User, cctv.Pwd))

	if err != nil {
		fmt.Printf("requsetError: %s\n", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("responseError: %s\n", err)
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		writeFile(cctv, imageOutPath, cctvCountMap, resp)
	case http.StatusUnauthorized:
		digestParts := util.GetDigestParts(resp)
		digestParts["uri"] = cctv.ImageUrlSuffix
		digestParts["method"] = req.Method
		digestParts["username"] = cctv.User
		digestParts["password"] = cctv.Pwd
		getDigestAuthImage(digestParts, cctvCountMap, cctv, url, imageOutPath)
	}
}

func getDigestAuthImage(digestParts map[string]string, cctvCountMap map[string]int,
	cctv conf.Cctv, url, imageOutPath string) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", util.GetDigestAuthrization(digestParts))

	if err != nil {
		fmt.Printf("requsetError: %s\n", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("responseError: %s\n", err)
		return
	}

	defer resp.Body.Close()

	writeFile(cctv, imageOutPath, cctvCountMap, resp)
}

func writeFile(cctv conf.Cctv, imageOutPath string, cctvCountMap map[string]int, resp *http.Response) {

	fileBytes, err := ioutil.ReadFile(imageOutPath + cctv.No + "_" + strconv.Itoa(cctvCountMap[cctv.No]) + ".jpg")

	if err != nil {
		fmt.Printf("readFileError: %s\n", err)
	}

	fileString := string(fileBytes)

	var fout *os.File
	switch len(fileString) == 0 {
	case true:
		fout, err = os.Create(imageOutPath + cctv.No + "_" + strconv.Itoa(cctvCountMap[cctv.No]) + ".jpg")
		//fmt.Println("創建檔案")
	case false:
		fout, err = os.OpenFile(imageOutPath+cctv.No+"_"+strconv.Itoa(cctvCountMap[cctv.No])+".jpg", os.O_RDWR|os.O_CREATE, 0666)
		//fmt.Println("開啟檔案")
	}

	if err != nil {
		fmt.Printf("fileOutError: %s\n", err)
	}

	defer fout.Close()

	switch {
	case cctv.ImageWidth == 0 || cctv.ImageHeight == 0:
		generateImageFile(cctv, resp, fout, cctvCountMap)
	case cctv.ImageWidth != 0 && cctv.ImageHeight != 0:
		generateCustomizeSizeImageFile(cctv, resp, fout, cctvCountMap)
	}
}

func generateImageFile(cctv conf.Cctv, resp *http.Response, fout *os.File, cctvCountMap map[string]int) {

	if _, err := io.Copy(fout, resp.Body); err != nil {
		fmt.Printf("copyError: %s\n", err)
	} else {
		fmt.Printf("寫入 <%s> 成功\n", util.CombineString(cctv.No, "_", strconv.Itoa(cctvCountMap[cctv.No])))

		if cctvCountMap[cctv.No] == 10 {
			cctvCountMap[cctv.No] = 1
		} else {
			cctvCountMap[cctv.No]++
		}
	}
}

func generateCustomizeSizeImageFile(cctv conf.Cctv, resp *http.Response, fout *os.File, cctvCountMap map[string]int) {

	bodys, err := ioutil.ReadAll(resp.Body)
	img, err := jpeg.Decode(bytes.NewReader(bodys))

	if err != nil {
		fmt.Printf("imgDecodeError: %s\n", err)
		return
	}

	newImage := resize.Resize(cctv.ImageWidth, cctv.ImageHeight, img, resize.MitchellNetravali)

	if err = jpeg.Encode(fout, newImage, &jpeg.Options{Quality: 40}); err != nil {
		fmt.Printf("imgEncodeError: %s\n", err)
	} else {
		fmt.Printf("寫入 <%s> 成功\n", util.CombineString(cctv.No, "_", strconv.Itoa(cctvCountMap[cctv.No])))

		if cctvCountMap[cctv.No] == 10 {
			cctvCountMap[cctv.No] = 1
		} else {
			cctvCountMap[cctv.No]++
		}
	}
}
