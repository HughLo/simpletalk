package xmlparser

/*
#cgo LDFLAGS: -lexpat
#include "expat.h"
#include <string.h>

*/
import "C"
import "errors"
import "unsafe"
import "fmt"

const (
    Ascii = "US-ASCII"
    UTF8  = "UTF-8"
    UTF16 = "UTF-16"
    ISO   = "ISO-8859-1"
)

type XPStatus int
const (
    InvalidStatus = iota
    Initialized
    Parsing
    Finished
    Suspended
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

func stringToXML_Char( from string ) *C.XML_Char {
    return (*C.XML_Char)(unsafe.Pointer(C.CString(from)))
}

func xml_boolToBool( from C.XML_Bool ) bool {
    ret := false
    if C.XML_TRUE == from {
        ret = true
    }

    return ret
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

func (self *XmlParser) Create( encoding string ) error {
    if !self.checkEncodingString( encoding ) {
        return fmt.Errorf("encoding %s is not supported", encoding)
    }

    x := stringToXML_Char( encoding )
    self.parserHandler = C.XML_ParserCreate(x)
    self.pinUserData = &userDataStructure{hooker:self.hooker}
    C.XML_SetUserData(self.parserHandler, unsafe.Pointer(self.pinUserData))
    return nil
}

func (self *XmlParser) Reset( encoding string ) error {
    if !self.checkEncodingString(encoding) {
        return fmt.Errorf("encoding %s is not supported", encoding)
    }

    x := stringToXML_Char(encoding)
    ret := xml_boolToBool(C.XML_ParserReset(self.parserHandler, x))
    if !ret {
        return self
    }

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
        return self
    }

    return nil
}

func (self *XmlParser) Stop(resumable bool) error {
    if self.parserHandler == nil {
        return errors.New("invalid parser handler")
    }

    c_resumable := 0
    if resumable {
        c_resumable = 1
    }
    retStatus := C.XML_StopParser(self.parserHandler, C.XML_Bool(C.uchar(c_resumable)))
    if retStatus != C.XML_STATUS_OK {
        return self
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

func (self *XmlParser) Status() int {
    statusRet := C.XML_ParsingStatus{}
    C.XML_GetParsingStatus(self.parserHandler, &statusRet)
    switch statusRet.parsing {
    case C.XML_INITIALIZED:
        return Initialized
    case C.XML_PARSING:
        return Parsing
    case C.XML_FINISHED:
        return Finished
    case C.XML_SUSPENDED:
        return Suspended
    }

    return InvalidStatus
}

func (self *XmlParser) StatusString() string {
    statusCode := self.Status()
    switch statusCode {
    case Initialized:
        return "initialized"
    case Parsing:
        return "parsing"
    case Finished:
        return "finished"
    case Suspended:
        return "suspended"
    }

    return "invalid status"
}

func (self *XmlParser) Error() string {
    errCode := C.XML_GetErrorCode(self.parserHandler)
    errStr := C.XML_ErrorString(errCode)
    return C.GoString((*C.char)(errStr))
}
