package xmlparser

/*
#cgo LDFLAGS: -lexpat
#include "expat.h"
#include <string.h>

*/
import "C"
import "errors"
import "unsafe"

const (
    Ascii = "US-ASCII"
    UTF8  = "UTF-8"
    UTF16 = "UTF-16"
    ISO   = "ISO-8859-1"
)

type userDataStructure struct {
    hooker *XmlParserHooker
    data interface{}
}

type XmlParser struct {
    parserHandler C.XML_Parser
    hooker *XmlParserHooker
    pinUserData *userDataStructure
}

func NewXmlParser() *XmlParser {
    parser := XmlParser{}
    parser.hooker = &XmlParserHooker{handlerMap:make(map[string]interface{})}
    return &parser
}

func (self *XmlParser) checkEncodingString( encoding string ) bool {
    switch {
    case encoding == Ascii, encoding == UTF8, encoding == UTF16, encoding == ISO:
        return true
    }
    return false
}

func (self *XmlParser) stringToXML_Char( from string ) *C.XML_Char {
    return (*C.XML_Char)(unsafe.Pointer(C.CString(from)))
}

func (self *XmlParser) Create( encoding string ) error {
    if !self.checkEncodingString( encoding ) {
        return errors.New( "encoding is not supported" )
    }

    x := self.stringToXML_Char( encoding )
    self.parserHandler = C.XML_ParserCreate(x)
    self.pinUserData = &userDataStructure{hooker:self.hooker}
    C.XML_SetUserData(self.parserHandler, unsafe.Pointer(self.pinUserData))
    return nil
}

func (self *XmlParser) Free() {
    if self.parserHandler != nil {
        C.XML_ParserFree( self.parserHandler )
    }
}

func (self *XmlParser) Parse( data string ) error {
    if self.parserHandler == nil {
        return errors.New( "invalid parser handler" )
    }

    cdata := C.CString(data)
    clen := C.strlen(cdata)

    retStatus := C.XML_Parse(self.parserHandler, cdata, C.int(clen), C.int(0))
    if retStatus != C.XML_STATUS_OK {
        return errors.New( "parser error" )
    }

    return nil
}

func (self *XmlParser) SetUserData(data interface{}) {
    self.pinUserData.data = data
    C.XML_SetUserData(self.parserHandler, unsafe.Pointer(self.pinUserData))
}

func (self *XmlParser) SetStartElementHandler(handler StartElementHandler) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = start_ele_handler
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)
}

func (self *XmlParser) SetEndElementHandler(handler EndElementHandler) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = end_ele_handler
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)
}
