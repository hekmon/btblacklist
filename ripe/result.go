package ripe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

/*
	Top level types
*/
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

/*
	intermediate types
*/
type doc struct {
	Strings answerStrings `json:"strs"`
}

func (d *doc) UnmarshalJSON(data []byte) (err error) {
	type alias doc
	tmp := struct {
		Doc struct {
			*alias
		} `json:"doc"`
	}{
		Doc: struct {
			*alias
		}{
			(*alias)(d),
		},
	}
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(&tmp)
}

type lst struct {
	Name    string         `json:"name"`
	Arrs    answerArrs     `json:"arrs"`
	Ints    answerIntegers `json:"ints"`
	Lists   []lst          `json:"lsts"`
	Strings answerStrings  `json:"strs"`
}

func (l *lst) UnmarshalJSON(data []byte) (err error) {
	type alias lst
	tmp := struct {
		List struct {
			*alias
		} `json:"lst"`
	}{
		List: struct {
			*alias
		}{
			(*alias)(l),
		},
	}
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(&tmp)
}

/*
	lower types
*/

type answerArrs map[string][]string

func (aa *answerArrs) UnmarshalJSON(data []byte) (err error) {
	// Unmarshall the fake JSON-XML payload
	var tmp []struct {
		Arr struct {
			Name   string `json:"name"`
			String struct {
				Value string `json:"value"`
			} `json:"str"`
		} `json:"arr"`
	}
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&tmp); err != nil {
		err = fmt.Errorf("can't unmarshall strings list into the tmp struct: %v", err)
		return
	}
	// Simplify it
	var exist bool
	*aa = make(map[string][]string, len(tmp))
	for _, value := range tmp {
		if _, exist = (*aa)[value.Arr.Name]; !exist {
			(*aa)[value.Arr.Name] = make([]string, 0, 1)
		}
		(*aa)[value.Arr.Name] = append((*aa)[value.Arr.Name], value.Arr.String.Value)
	}
	return
}

type answerIntegers map[string][]int

func (ai *answerIntegers) UnmarshalJSON(data []byte) (err error) {
	// Unmarshall the fake JSON-XML payload
	var tmp []struct {
		Integer struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"int"`
	}
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&tmp); err != nil {
		err = fmt.Errorf("can't unmarshall strings list into the tmp struct: %v", err)
		return
	}
	// Simplify it
	var (
		exist  bool
		tmpInt int
	)
	*ai = make(map[string][]int, len(tmp))
	for _, value := range tmp {
		if _, exist = (*ai)[value.Integer.Name]; !exist {
			(*ai)[value.Integer.Name] = make([]int, 0, 1)
		}
		if tmpInt, err = strconv.Atoi(value.Integer.Value); err != nil {
			err = fmt.Errorf("can't convert integer str as int: %v", err)
			return
		}
		(*ai)[value.Integer.Name] = append((*ai)[value.Integer.Name], tmpInt)
	}
	return
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
