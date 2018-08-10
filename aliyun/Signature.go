package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

const (
	access_key        = "access_key"
	access_key_secret = "access_key_secret"
	ecs_url           = "https://ecs.aliyuncs.com/"
)

func GetUrlFormedMap(source map[string]string) (urlEncoded string) {
	urlEncoder := url.Values{}
	for key, value := range source {
		urlEncoder.Add(key, value)
	}
	urlEncoded = urlEncoder.Encode()
	return
}

func Sign(stringToSign, secretSuffix string) string {
	secret := access_key_secret + secretSuffix
	return ShaHmac1(stringToSign, secret)
}

func ShaHmac1(source, secret string) string{
	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(source))
	signedBytes := mac.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}

func GetSign(signParams map[string]string)  (signature string) {
	var paramKeySlice []string
	for key := range signParams {
		paramKeySlice = append(paramKeySlice, key)
	}
	sort.Strings(paramKeySlice)
	stringToSign := GetUrlFormedMap(signParams)
	stringToSign = strings.Replace(stringToSign, "+", "%20", -1)
	stringToSign = strings.Replace(stringToSign, "*", "%2A", -1)
	stringToSign = strings.Replace(stringToSign, "%7E", "~", -1)
	stringToSign = url.QueryEscape(stringToSign)
	stringToSign = "GET" + "&%2F&" + stringToSign

	signature = Sign(stringToSign, "&")
	return
}

func GetUUIDV4() (uuidHex string) {
	uuidV4, _ := uuid.NewV4()
	uuidHex = hex.EncodeToString(uuidV4.Bytes())
	return
}

var LoadLocationFromTZData func(name string, data []byte) (*time.Location, error) = nil
var TZData []byte = nil

func GetGMTLocation() (*time.Location, error) {
	if LoadLocationFromTZData != nil && TZData != nil {
		return LoadLocationFromTZData("GMT", TZData)
	} else {
		return time.LoadLocation("GMT")
	}
}
func GetTimeInFormatISO8601() (timeStr string) {
	gmt, err := GetGMTLocation()

	if err != nil {
		panic(err)
	}
	return time.Now().In(gmt).Format("2006-01-02T15:04:05Z")
}



func EcsAction(actionname string, paras1 map[string]string) string {
	paras := map[string]string{
		"SignatureVersion": "1.0",
		"Format":           "JSON",
		"Timestamp":        GetTimeInFormatISO8601(),
		"AccessKeyId":      access_key,
		"SignatureMethod":  "HMAC-SHA1",
		"Version":          "2014-05-26",
		"Action":           actionname,
		"SignatureNonce":   GetUUIDV4(),
	}
	for k, v := range paras1 {
		paras[k] = v
	}
	paras["Signature"] = GetSign(paras)
	queryString := GetUrlFormedMap(paras)
	req, _ := http.NewRequest("GET", ecs_url+"?"+queryString,strings.NewReader(""))
	client := http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	m1 := make(map[string]interface{})
	json.Unmarshal(body, &m1)
	fJson, _ := json.MarshalIndent(m1, "", "  ")
	return string(fJson)
}

func main() {

	m2 := map[string]string{
		"Generation": "ecs-1",
		"RegionId":   "cn-beijing",
	}
	res := EcsAction("DescribeInstanceTypeFamilies", m2)
	fmt.Printf("%s", res)

}

