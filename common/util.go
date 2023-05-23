package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Float32ToString(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

func StringToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int(i)
}

func StringToInt32(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}

func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func StringToFloat32(s string) float32 {
	f, _ := strconv.ParseFloat(s, 32)
	return float32(f)
}

func StringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func ByteArrayToHexString(byteArray []byte) string {
	var buffer bytes.Buffer
	for _, b := range byteArray {
		//builder.append(Integer.toString((arrayBytes[i] & 0xff) + 0x100, 16).substring(1));
		buffer.WriteString(strconv.FormatInt((int64(b)&0xff)+0x100, 16)[1:])
	}
	return buffer.String()
}

/**
MapGet : Tries to get the key from the map m, if not found it returns defaultValue
*/
func MapGet(m map[string]string, key, defaultValue string) string {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

func RoundDecimalDigitsInStr(f float64, d int) string {
	return strconv.FormatFloat(f, 'f', d, 64)
}

func GetDateFormatFromString(year string, month string) (time.Time, error) {
	value := year + "-" + month
	layout := "2006-01"
	t, err := time.Parse(layout, value)
	if err != nil {
		fmt.Println(err)
		//Log into logs
	}
	return t, nil
}

func DecodeBase64(encodedStr string) string {
	decodedStr, _ := base64.StdEncoding.DecodeString(encodedStr)
	return string(decodedStr)
}
func TimeDiffFromNow(dateTime time.Time, nowTime time.Time) string {
	resultDuration := time.Until(dateTime)
	resultHours := resultDuration.Hours()

	noOfDays := int64(resultHours / 24)

	return Int64ToString(noOfDays)

}

func GetTimeFromDateTimeWithSeconds(dateTime time.Time) string {
	return dateTime.Format("15:04:05")
}

func GetTimeFromDateTimeWithoutSeconds(dateTime time.Time) string {
	return dateTime.Format("15:04")
}

func GetDateFromDateTime(dateTime time.Time) string {
	return dateTime.Format("2006-01-02")
}

//This function should be used for incoming request date conversion
func GetDateFormString(dateStr string, format *string) (time.Time, error) {
	if format == nil || *format == "" {
		return time.Parse(time.RFC3339Nano, dateStr)
	}
	return time.Parse(*format, dateStr)
}

func Guid() string {
	return uuid.NewV4().String()
}

func AddGetParamsToUrl(urlPath string, queryParams url.Values) (string, error) {
	u, err := url.Parse(urlPath)
	if err == nil {
		finalQueryParams, err := url.ParseQuery(u.RawQuery)
		if err == nil {
			for k, v := range queryParams {
				if ov, ok := finalQueryParams[k]; ok {
					ov = append(ov, v...)
					v = ov
				}
				finalQueryParams[k] = v
			}
			u.RawQuery = finalQueryParams.Encode()
		}
		return u.String(), nil
	}
	return "", err
}

func StripOffIsdFromMobile(mobile string) string {
	if len(mobile) > 10 {
		if mobile[0:2] == "91" {
			return mobile[2:]
		}
	}
	return mobile
}

func CompareMobileNumber(mobile1, mobile2 string) bool {
	mobile1 = StripOffIsdFromMobile(mobile1)
	mobile2 = StripOffIsdFromMobile(mobile2)
	return mobile1 == mobile2
}

func PtrString(s string) *string {
	return &s
}

func StrFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func IntFromPtr(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
func Int64FromPtr(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// This returns the caller of the function that called it
func GetCallerFullFuncName(skip int) string {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	skip = skip + 3
	n := runtime.Callers(skip, fpcs)
	if n == 0 {
		return "n/a" // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}
	// return its name
	return fun.Name()
}

func AmtFromX100(x100 int64) float64 {
	return float64(x100) / float64(100)
}

func ObjToJsonStr(obj interface{}) string {
	if b, err := json.Marshal(obj); err != nil {
		return ""
	} else {
		return string(b)
	}

}

//Multipart file upload
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func EscapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func UrlFromPathParams(url string, pathParams map[string]string) string {
	for k, v := range pathParams {
		url = strings.Replace(url, "{"+k+"}", v, -1)
	}
	return url
}

func TruncateDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func TruncateMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, t.Location())
}

func TruncateYear(t time.Time) time.Time {
	return time.Date(t.Year(), 0, 0, 0, 0, 0, 0, t.Location())
}

func CalculateEmi(totalCostX100 int64, ratePerAnnum float64, months int64) float64 {
	roiPerMonth := ratePerAnnum / 1200   //R is ROI per month
	intRatePlus := roiPerMonth + 1       //1+R
	var emiCalculatingFactor float64 = 1 //(1+R)^N
	var emi float64                      //Monthly instalment
	if roiPerMonth == 0 {
		emi = float64(totalCostX100 / months)
	} else {
		for i := int64(1); i <= months; i++ {
			emiCalculatingFactor = emiCalculatingFactor * intRatePlus
		}
		emi = float64(totalCostX100) * roiPerMonth * emiCalculatingFactor / (emiCalculatingFactor - 1)
	}
	return emi
}

//VersionCheck - returns +ve if a>b else -ve if b>a else if equal then 0.
//versions should be of form x.x.x and so on
func VersionCheck(a, b string) int {
	aSplits := strings.Split(a, ".")
	bSplits := strings.Split(b, ".")
	la := len(aSplits)
	lb := len(bSplits)
	l := la
	if l < lb {
		l = lb
	}
	var iA, iB int64
	for i := 0; i < l; i++ {
		if la <= i {
			iA = 0
		} else {
			iA, _ = strconv.ParseInt(aSplits[i], 10, 64)
		}
		if lb <= i {
			iB = 0
		} else {
			iB, _ = strconv.ParseInt(bSplits[i], 10, 64)
		}
		if iA > iB {
			return i + 1
		} else if iB > iA {
			return -(i + 1)
		}
	}
	return 0
}
