package xmlparser

import "testing"
import "log"

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

var xmlData string = `
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
var xmlDataForUserData string = "<test></test>"

//TestParseWhole tests the entire xml parsing ability
func TestParseWhole( t *testing.T ) {
    testHassle := new(xmlParserTestHassle)
    testHassle.SetupHassle()
    defer testHassle.DestroyHassle()

    startEleHandler := func (userData interface{}, name string, attrs map[string]string) {
        log.Print("<"+name+">")
        for key, val := range attrs {
            log.Print( key + ":" + val)
        }
    }
    endEleHandler := func (userData interface{}, name string) {
        log.Printf( "</%s>\n", name )
    }

    testHassle.parser.SetStartElementHandler(startEleHandler)
    testHassle.parser.SetEndElementHandler(endEleHandler)
    testHassle.parser.Parse(xmlData)
}

//TestUserData tests user data can work fine
func TestUserData( t *testing.T ) {
    testHassle := new(xmlParserTestHassle)
    testHassle.SetupHassle()
    defer testHassle.DestroyHassle()

    var result string
    startEleHandler := func(userData interface{}, name string, attrs map[string]string) {
        if userData == nil {
            result = "empty"
            return
        }

        data, ok := userData.(*string)
        if !ok {
            result = "mismatch"
            return
        }

        result = *data
    }

    endEleHandler := func(userData interface{}, name string) {
        if userData == nil {
            result = "empty"
            return
        }

        data, ok := userData.(*string)
        if !ok {
            result = "mismatch"
            return
        }

        result = *data
    }

    testUserData := "test user data"
    testHassle.parser.SetUserData(&testUserData)
    testHassle.parser.SetStartElementHandler(startEleHandler)
    testHassle.parser.SetEndElementHandler(endEleHandler)

    testHassle.parser.Parse("<test>")
    if result != "test user data" {
        t.Error( "not equal" )
    }

    testHassle.parser.SetUserData(nil)
    testHassle.parser.Parse( "<test>" )
    if result != "empty" {
        t.Errorf( "not equal:%s", result )
    }
    testHassle.parser.Parse("</test>")
    if result != "empty" {
        t.Errorf( "not equal:%s", result )
    }
    testHassle.parser.Parse("</test>")
    if result != "empty" {
        t.Error( "not equal" )
    }

    err := testHassle.parser.Parse("<test>")
    if err != nil {
        log.Println( err )
    }
}

func TestStatus(t *testing.T) {
    testHassle := new(xmlParserTestHassle)
    testHassle.SetupHassle()
    defer testHassle.DestroyHassle()

    parseStatus := testHassle.parser.Status()
    if parseStatus != Initialized {
        t.Errorf("incorrect status: %d\n", parseStatus)
    }

    startEleHandler := func(ud interface{}, name string, attrs map[string]string) {
    }

    endEleHandler := func(ud interface{}, name string) {
        /*
        parser, ok := ud.(*XmlParser)
        if !ok {
            t.Errorf("user data type mismatch\n")
        }

        parser.Stop(true)
        */
    }

    testHassle.parser.SetUserData(testHassle.parser)
    testHassle.parser.SetStartElementHandler(startEleHandler)
    testHassle.parser.SetEndElementHandler(endEleHandler)
    testHassle.parser.Parse("<test>")
    parseStatus = testHassle.parser.Status()
    if parseStatus != Parsing {
        t.Errorf("incorrect status: %s\n", testHassle.parser.StatusString())
    }

    testHassle.parser.Parse("</test>")
    parseStatus = testHassle.parser.Status()
    if parseStatus != Finished {
        t.Errorf("incorrect status: %s\n", testHassle.parser.StatusString())
    }
}

func TestNewXmlParser(t *testing.T) {
    parser := NewXmlParser()
    if parser == nil {
        t.FailNow()
    }
}
