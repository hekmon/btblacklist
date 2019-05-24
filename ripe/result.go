package ripe

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
		Strings []answerString `json:"strs"`
	} `json:"doc"`
}

type lst struct {
	List struct {
		Name    string          `json:"name"`
		Arrs    []answerArr     `json:"arrs"`
		Ints    []answerInteger `json:"ints"`
		Lists   []lst           `json:"lsts"`
		Strings []answerString  `json:"strs"`
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

type answerString struct {
	String struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"str"`
}
