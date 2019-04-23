package main

import (
	"encoding/json"
	"fmt"
	"github.com/BillSJC/go-curl"
)

type transData struct {
	Err       error
	ErrorCode int
	Msg       string
	Data      interface{}
}

func newTransData(err error, errorCode int, msg string, data interface{}) *transData {
	return &transData{
		Err:       err,
		ErrorCode: errorCode,
		Msg:       msg,
		Data:      data,
	}
}

func (s *transData) GetData() interface{} {
	return s.Data
}

/* const  */
const (
	constRequestScheduleURL   = "url to get student info"
	constRequestPersonInfoURL = "url to get person info"
	accessToken               = "set accessToken here"
)

func main() {

}

func getData() {

	personChan := make(chan *transData)
	go sendRequest(constRequestPersonInfoURL, accessToken, personChan, new(person))

	scoreChan := make(chan *transData)
	go sendRequest(constRequestScheduleURL, accessToken, scoreChan, new(grades))

	//get data
	personDataRaw := <-personChan
	scoreDataRaw := <-scoreChan

	closeAllChan(personChan, scoreChan)

	if personDataRaw.Err != nil || scoreDataRaw.Err != nil {
		panic("something error")
		return
	}

	personData := personDataRaw.Data.(*person)
	gradeData := scoreDataRaw.Data.(*grades)

	fmt.Println(personData)
	fmt.Println(gradeData)
}

func closeAllChan(inputs ...chan *transData) {
	for _, v := range inputs {
		close(v)
	}
}

//send request by chan
func sendRequest(reqURL string, token string, dataChan chan *transData, data interface{}) {
	resp, err := curl.NewRequest().SetUrl(reqURL).SetHeaders(map[string]string{"authorization": "token " + token}).Get()
	if err != nil || resp == nil {
		dataChan <- newTransData(err, 500, "", data)
		return
	}
	if !resp.IsOk() {
		dataChan <- newTransData(fmt.Errorf("code:%d |||data:%s", resp.Raw.StatusCode, resp.Body), 500, "", data)
		return
	}
	err = json.Unmarshal([]byte(resp.Body), &data)
	if err != nil {
		dataChan <- newTransData(err, 500, "", data)
		return
	}
	dataChan <- newTransData(nil, 0, "success", data)
}

/*	define struct here	*/

type grades struct {
	Error int
	Msg   string
	Data  []*gradeData
}
type gradeData struct {
	COURSE     string
	COURSECODE string
	COURSETYPE string
	CREDIT     string
	SCHOOLYEAR string
	SCORE      string
	SCORE_PS   string
	SCORE_QM   string
	SCORE_QZ   string
	SCORE_SY   string
	SELECTCODE string
	SEMESTER   string
	SEORE_BK   string
	STAFFID    string
	STAFFNAME  string
}

type person struct {
	Error   int
	Msg     string
	IsCache int
	Data    *person_info
}
type person_info struct {
	STAFFID   string
	STAFFNAME string
	STAFFTYPE string
	STAFSTATE string
	UNITCODE  string
}
