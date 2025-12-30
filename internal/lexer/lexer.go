package lexer

// Token represents a lexical token.
type Token struct {
    Type  string
    Value string
}

// Lexer reads characters and produces tokens.
type Lexer struct {
    src string
    pos int
}

func New(src string) *Lexer {
    return &Lexer{src: src}
}

func (l *Lexer) Lex() []Token {
    tokens := []Token{}

    for l.pos < len(l.src) {
        ch := l.src[l.pos]

        switch {
        case isLetter(ch):
            ident := l.readIdent()
            tokens = append(tokens, Token{Type: "IDENT", Value: ident})

        case isDigit(ch):
            num := l.readNumber()
            tokens = append(tokens, Token{Type: "NUMBER", Value: num})

        case ch == '"':
            str := l.readString()
            tokens = append(tokens, Token{Type: "STRING", Value: str})

        default:
            tokens = append(tokens, Token{Type: string(ch), Value: string(ch)})
            l.pos++
        }
    }

    return tokens
}

func isLetter(ch byte) bool {
    return (ch >= 'a' && ch <= 'z') ||
        (ch >= 'A' && ch <= 'Z') ||
        ch == '_'
}

func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
}

func (l *Lexer) readIdent() string {
    start := l.pos
    for l.pos < len(l.src) && isLetter(l.src[l.pos]) {
        l.pos++
    }
    return l.src[start:l.pos]
}

func (l *Lexer) readNumber() string {
    start := l.pos
    for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
        l.pos++
    }
    return l.src[start:l.pos]
}

func (l *Lexer) readString() string {
    l.pos++ // skip opening quote
    start := l.pos
    for l.pos < len(l.src) && l.src[l.pos] != '"' {
        l.pos++
    }
    str := l.src[start:l.pos]
    l.pos++ // skip closing quote
    return str
}
