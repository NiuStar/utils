package utils

import (
	//"fmt"
	"crypto/md5"
	"crypto/sha1"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"runtime"

	"time"
)

// 获取文件大小的接口
type Sizer interface {
	Size() int64
}

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

//var coder = base64.NewEncoding(base64Table)

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

/*********
reqString = Substr(reqString, 0, len(reqString)-1)
*********/
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {

		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

/*********
ltime : 秒数
dataformat : "2006/01/02 - 03:04:05"
**********/

func FormatTime(ltime int64, dataformat string) string {
	tm := time.Unix(ltime, 0)
	return tm.Format(dataformat)
}

func FormatTimeAll(ltime int64) string {
	tm := time.Unix(ltime, 0)
	return tm.Format("2006-01-02 15:04:05")
}

func ReadFile(path string) string {
	path = GetCurrPath() + path
	fmt.Println("path:  " , path)
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}


func ReadFileFullPath(path string) string {


	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}


func Base64Encode(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
	//return coder.EncodeToString([]byte(str))
}

func Base64Decode(str string) (string,bool) {
	str = strings.Replace(str, " ", "+", -1)
	data, err := b64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
		return "",false
	}
	return string(data),true
}

func FormatImageBase64FormURL(data string) string {
	jpg_head := "data:image/jpeg;base64,"
	png_head := "data:image/png;base64,"
	imgData := ""
	if strings.Index(data,jpg_head) == 0 {
		imgData = Substr(data, len(jpg_head), len(data)-len(jpg_head))
	} else if strings.Index(data,png_head) == 0 {
		imgData = Substr(data, len(png_head), len(data)-len(png_head))
	} else {
		imgData = data
	}
	return imgData
}

func WriteToFile(path string, data string) error {
	return ioutil.WriteFile(path, []byte(data), 0666) //写入文件(字节数组)
}

func GetNetPath(uri string,path string) string {
	s := os.Args[0]
	fmt.Println("os.Args[0] : ",s)
	var index int = 0
	if runtime.GOOS == "windows" {
		index = strings.LastIndex(s,"\\")



	} else {
		index = strings.LastIndex(s,"/")
	}

	fmt.Println("index :",index)
	s = s[0:index + 1]

	if strings.LastIndex(uri,"/") != len(uri) - 1 {
		uri += "/"
	}

	return uri + path[len(s):len(path)]

}


func Mkdir(path string) bool {
	err := os.MkdirAll(path, 0666)
	if err != nil {
		fmt.Printf("%s", err)
		return false
	} else {
		fmt.Print("Create Directory OK!")
	}
	return true
}
func SaveUploadFile(file multipart.File, filename string) {
	path := "./html/upload/" + filename
	nFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
		os.Exit(-1)
	}
	defer nFile.Close()
	io.Copy(nFile, file)
}
func GetFileLenght(file multipart.File) (int64, bool) {
	if sizeInterface, ok := file.(Sizer); ok {
		if sizeInterface != nil {
			return sizeInterface.Size(), true
		} else {
			return 0, false
		}
	} else {
		return 0, false
	}
}

var config map[string]interface{} = make(map[string]interface{})

func ReadUtils() {
	data := ReadFile("utils.json")
	err := json.Unmarshal([]byte(data), &config)
	if err != nil {
		fmt.Println("utils文件错误")
		panic(err)
		return
	}
}

func GetUtilsValue(key string) string {
	return config[key].(string)
}

func MD5(str string) string {

	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串为 sharejs.com
	return hex.EncodeToString(h.Sum(nil))

}

func SHA1(str string) string {

	h := sha1.New()
	h.Write([]byte(str)) // 需要加密的字符串为 sharejs.com
	return hex.EncodeToString(h.Sum(nil))

}



func Int2Unicode(b int) string {
	s := fmt.Sprintf("%#U", b)
	arr := strings.Split(s,`'`)
	if len(arr) == 3 {
		return arr[1]
	} else {
		//fmt.Println(s)
		return s
	}

}

func GET(uri string) string {
	var client = &http.Client{}
	//向服务端发送get请求

	request, _ := http.NewRequest("GET", uri, nil)

	response, err := client.Do(request)

	if err != nil {
		panic(err)
		return ""
	}


	if response.StatusCode == 200 {

		body, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()
		return string(body)

	}
	response.Body.Close()
	return ""
}

func POST(uri string, datas map[string]string) string {
	postValues := url.Values{}
	for key, value := range datas {
		postValues.Add(key, value)
	}
	body := ioutil.NopCloser(strings.NewReader(postValues.Encode())) //把form数据编下码
	var client = &http.Client{}
	//向服务端发送get请求

	req, _ := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	resp, err := client.Do(req)

	if err != nil {
		// handle error
		panic(err)
		return "{}"
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
		// handle error
	}
	resp.Body.Close()
	return string(data)

}
