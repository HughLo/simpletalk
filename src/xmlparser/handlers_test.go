package xmlparser

import (
    "testing"
)

func TestHook( t *testing.T ) {
    parser := NewXmlParser()
    parserHandler := new(parser_handler)
    hooker := new(XmlParserHooker)
    hooker.Hook( parser, parserHandler)
}
