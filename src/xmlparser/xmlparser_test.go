package xmlparser

import "testing"
import "log"

func HandleStartEle(userData interface{}, name string, attrs map[string]string) {
    log.Print("<"+name+">")
    for key, val := range attrs {
        log.Print( key + ":" + val)
    }
}

func HandleEndEle(userData interface{}, name string) {
    _, ok := userData.(*testing.T)
    if !ok {
        log.Fatalln( "HandleEndEle user data is not of type userDataStructure" )
    }

    log.Printf( "/name: %s\n", name )
}

type xmlParserTestHassle struct {
    parser *XmlParser
    data []*string
}

func (self *xmlParserTestHassle) SetupHassle() {
    self.parser = NewXmlParser()
    self.parser.Create(UTF8)
    self.data = make([]*string, 10)
}

func (self *xmlParserTestHassle) DestroyHassle() {
    if self.parser != nil {
        self.parser.Free()
    }
}


func TestParse( t *testing.T ) {
    testHassle := new(xmlParserTestHassle)
    testHassle.SetupHassle()
    defer testHassle.DestroyHassle()

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
    testHassle.parser.SetStartElementHandler(HandleStartEle)
    testHassle.parser.SetEndElementHandler(HandleEndEle)
    testHassle.parser.SetUserData(t)
    testHassle.parser.Parse(data)
}

func TestNewXmlParser(t *testing.T) {
    parser := NewXmlParser()
    if parser == nil {
        t.FailNow()
    }
}
