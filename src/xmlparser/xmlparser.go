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

func (self *XmlParser) callHook(handlerName string, handler interface{}) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = handlerName
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)

}

//Create create the parser to handle the xml stream encoded by encoding
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

//Reset reset the parser
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

//Free free the parser
func (self *XmlParser) Free() {
    if self.parserHandler != nil {
        C.XML_ParserFree( self.parserHandler )
    }
}

//Parse parser the content specifed by data
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

//Stop stop the parser
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

//SetUserData set user data which will be passed as the first argument of callbacks
func (self *XmlParser) SetUserData(data interface{}) {
    self.pinUserData.data = data
    C.XML_SetUserData(self.parserHandler, unsafe.Pointer(self.pinUserData))
}

//SetStartElementHandler set start element handler
func (self *XmlParser) SetStartElementHandler(handler StartElementHandler) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = start_ele_handler
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)
}

//SetEndElementHandler set end element handler
func (self *XmlParser) SetEndElementHandler(handler EndElementHandler) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = end_ele_handler
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)
}

//SetCharacterDataHandler set handler of characters
func (self *XmlParser) SetCharacterDataHandler(handler CharacterDataHandler) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = character_data_handler
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)
}

//SetProcessingInstHandler set processing instruction handler
func (self *XmlParser) SetProcessingInstHandler(handler PIHandler) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = processing_inst_handler
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)
}

//SetCommentHandler set processing comment handler
func (self *XmlParser) SetCommentHandler(handler CommentHandler) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = comment_handler
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)
}

//SetStartCDataSectionHandler set handler of start CDATA section
func (self *XmlParser) SetStartCDataSectionHandler(handler StartCDataSectionHandler) error {
    var handlerData interface{} = handler
    if handler == nil {
        null_handler.name = start_cdata_section_handler
        handlerData = null_handler
    }
    return self.hooker.Hook(self, handlerData)
}

//SetEndCDataSectionHandler set handler of end CDATA section
func (self *XmlParser) SetEndCDataSectionHandler(handler EndCDataSectionHandler) error {
    return self.callHook(end_cdata_section_handler, handler)
}

//SetDefaultHandler set default handler
func (self *XmlParser) SetDefaultHandler(handler DefaultHandler) error {
    return self.callHook(default_handler, handler)
}

//SetHandlers set all handlers at one time
func (self *XmlParser) SetHandlers(handlers *map[string]interface{}) error {
    if handlers == nil {
        return errors.New("invalid parameters")
    }

    for k, v := range *handlers {
        switch k {
        case start_ele_handler:
            var trueHandler StartElementHandler
            if v != nil {
                trueHandler, ok := v.(StartElementHandler)
                if !ok || trueHandler == nil {
                    return errors.New("handler type mismatch")
                }
            }
            self.SetStartElementHandler(trueHandler)
        case end_ele_handler:
            var trueHandler EndElementHandler
            if v!= nil {
                trueHandler, ok := v.(EndElementHandler)
                if !ok || trueHandler == nil {
                    return errors.New("handler type mismatch")
                }
            }
            self.SetEndElementHandler(trueHandler)
        case character_data_handler:
            var trueHandler CharacterDataHandler
            if v!= nil {
                trueHandler, ok := v.(CharacterDataHandler)
                if !ok || trueHandler == nil {
                    return errors.New("handler type mismatch")
                }
            }
            self.SetCharacterDataHandler(trueHandler)
        case processing_inst_handler:
            var trueHandler PIHandler
            if v!= nil {
                trueHandler, ok := v.(PIHandler)
                if !ok || trueHandler == nil {
                    return errors.New("handler type mismatch")
                }
            }
            self.SetProcessingInstHandler(trueHandler)
        case comment_handler:
            var trueHandler CommentHandler
            if v!= nil {
                trueHandler, ok := v.(CommentHandler)
                if !ok || trueHandler == nil {
                    return errors.New("handler type mismatch")
                }
            }
            self.SetCommentHandler(trueHandler)
        case start_cdata_section_handler:
            var trueHandler StartCDataSectionHandler
            if v!= nil {
                trueHandler, ok := v.(StartCDataSectionHandler)
                if !ok || trueHandler == nil {
                    return errors.New("handler type mismatch")
                }
            }
            self.SetStartCDataSectionHandler(trueHandler)
        case end_cdata_section_handler:
            var trueHandler EndCDataSectionHandler
            if v!= nil {
                trueHandler, ok := v.(EndCDataSectionHandler)
                if !ok || trueHandler == nil {
                    return errors.New("handler type mismatch")
                }
            }
            self.SetEndCDataSectionHandler(trueHandler)
        case default_handler:
            var trueHandler DefaultHandler
            if v!= nil {
                trueHandler, ok := v.(DefaultHandler)
                if !ok || trueHandler == nil {
                    return errors.New("handler type mismatch")
                }
            }
            self.SetDefaultHandler(trueHandler)
        default:
            return errors.New("unsupported handler type")
        }
    }

    return nil
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
