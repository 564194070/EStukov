package util

import (
	"encoding/json"
	"log"
)

type JsonMessage struct {
	ResCode int	`json:"code"`
	Message string `json:"message"`
	Extend interface{} `json:"data"`
}

func NewReqMessage(resCode int, message string, extend interface{}) *JsonMessage {
	return & JsonMessage{
		ResCode: resCode,
		Message: message,
		Extend: extend,
	}
}
func (j *JsonMessage)JsonToByte() []byte {
	res, err := json.Marshal(j)
	if err != nil {
		log.Fatalln("Failed J tyo B, err", err)
	}
	return res
}

func (j *JsonMessage)JsonToString() string {
	res, err := json.Marshal(j)
	if err != nil {
		log.Fatalln("Failed J tyo B, err", err)
	}
	return string(res)
}
