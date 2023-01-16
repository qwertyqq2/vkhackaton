package htmlparser

type token struct {
	n      string
	data   string
	closed bool
	parent *token
}

type tokensList struct {
	t []*token
}

func NewTokensList() *tokensList {
	return &tokensList{
		t: make([]*token, 0),
	}
}

func (l *tokensList) addToken(n string) {
	l.t = append(l.t, &token{
		n:      n,
		closed: false,
		data:   "",
	})
	if len(l.t) > 1 {
		for i := len(l.t) - 2; i >= 0; i-- {
			if l.t[i].closed == false {
				l.t[len(l.t)-1].parent = l.t[i]
				break
			}
		}
	} else {
		l.t[len(l.t)-1].parent = nil
	}
}

func (l *tokensList) closeToken(n string) {
	for i := len(l.t) - 1; i >= 0; i-- {
		if l.t[i].closed == false {
			l.t[i].closed = true
			break
		}
	}
}

func (l *tokensList) text(t string) {
	f := false
	for _, r := range t {
		if r != 32 && r != 10 {
			f = true
		}
	}
	if !f {
		return
	}
	for i := len(l.t) - 1; i >= 0; i-- {
		if l.t[i].closed == false {
			l.t[i].data += t
			break
		}
	}
}

func addChilders(t *token, str string, tokens []*token) string {
	for _, tn := range tokens {
		if tn.parent == t {
			str += addChilders(tn, tn.data, tokens) + "\n"
		}
	}
	return str
}

func (l *tokensList) body() string {
	var bodyT *token
	bodyData := ""
	for _, tn := range l.t {
		if tn.n == "<body>" {
			bodyT = tn
		}
	}
	return addChilders(bodyT, bodyData, l.t)
}

func (l *tokensList) head() string {
	var headT *token
	headData := ""
	for _, tn := range l.t {
		if tn.n == "<head>" {
			headT = tn
		}
	}
	return addChilders(headT, headData, l.t)
}
