package ripe

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type answer struct {
	Result answerResult `json:"result"`
	Lists  []lst        `json:"lsts"`
}

type answerResult struct {
	Name     string `json:"name"`
	NumFound int    `json:"numFound"`
	Start    int    `json:"start"`
	Docs     []doc  `json:"docs"`
}

type doc struct {
	Doc struct {
		Strings answerStrings `json:"strs"`
	} `json:"doc"`
}

type lst struct {
	List struct {
		Name    string          `json:"name"`
		Arrs    []answerArr     `json:"arrs"`
		Ints    []answerInteger `json:"ints"`
		Lists   []lst           `json:"lsts"`
		Strings answerStrings   `json:"strs"`
	} `json:"lst"`
}

/*
	base units
*/

type answerArr struct {
	Arr struct {
		Name   string `json:"name"`
		String struct {
			Value string `json:"value"`
		} `json:"str"`
	} `json:"arr"`
}

type answerInteger struct {
	Integer struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"int"`
}

type answerStrings map[string][]string

func (as *answerStrings) UnmarshalJSON(data []byte) (err error) {
	// Unmarshall the fake JSON-XML payload
	var tmp []struct {
		String struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"str"`
	}
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&tmp); err != nil {
		err = fmt.Errorf("can't unmarshall strings list into the tmp struct: %v", err)
		return
	}
	// Simplify it
	var exist bool
	*as = make(map[string][]string, len(tmp))
	for _, value := range tmp {
		if _, exist = (*as)[value.String.Name]; !exist {
			(*as)[value.String.Name] = make([]string, 0, 1)
		}
		(*as)[value.String.Name] = append((*as)[value.String.Name], value.String.Value)
	}
	return
}
