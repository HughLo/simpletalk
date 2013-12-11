package xmlparser

import "testing"
import "log"

type parser_handler struct {
}

func (self parser_handler) HandleStartEle(userData interface{}, name string, attrs map[string]string) {
    log.Print("<"+name+">")
    for key, val := range attrs {
        log.Print( key + ":" + val)
    }
}

func (self parser_handler) HandleEndEle(userData interface{}, name string) {
    log.Print( "</"+name+">" )
}

func TestParse( t *testing.T ) {
    data := `
		<Person type="test">
			<FullName>Grace R. Emlin</FullName>
			<Company>Example Inc.</Company>
			<Email where="home">
				<Addr>gre@example.com</Addr>
			</Email>
			<Email where='work'>
				<Addr>gre@work.com</Addr>
			</Email>
			<Group>
				<Value>Friends</Value>
				<Value>Squash</Value>
			</Group>
			<City>Hanga Roa</City>
			<State>Easter Island</State>
		</Person>
	`
    parser := NewXmlParser()
    handler := new(parser_handler)
    parser.Create(UTF8)
    parser.SetStartElementHandler(handler)
    parser.SetEndElementHandler(handler)
    parser.Parse(data)
    parser.Free()
}

func TestNewXmlParser(t *testing.T) {
    parser := NewXmlParser()
    if parser == nil {
        t.FailNow()
    }
}
