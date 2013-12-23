package xmlparser

/*
#cgo LDFLAGS: -lexpat
#include <expat.h>
extern int getExpatArrayLen( char** data );
extern void hookStartElementHandler(XML_Parser parser);
extern void unhookStartElementHandler(XML_Parser parser);
extern void hookEndElementHandler(XML_Parser parser);
extern void unhookEndElementHandler(XML_Parser parser);
extern void hookCharacterDataHandler(XML_Parser parser);
extern void unhookCharacterDataHandler(XML_Parser parser);
extern void hookPIHandler(XML_Parser parser);
extern void unhookPIHandler(XML_Parser parser);
extern void hookCommentHandler(XML_Parser parser);
extern void unhookCommentHandler(XML_Parser parser);
extern void hookStartCDataSectionHandler(XML_Parser parser);
extern void unhookStartCDataSectionHandler(XML_Parser parser);
extern void hookEndCDataSectionHandler(XML_Parser parser);
extern void unhookEndCDataSectionHandler(XML_Parser parser);
extern void hookDefaultHandler(XML_Parser parser);
extern void unhookDefaultHandler(XML_Parser parser);
*/
import "C"
import (
    "errors"
    "unsafe"
    "log"
    "reflect"
    _ "fmt"
)

const (
    start_ele_handler = "start_ele_handler"
    end_ele_handler = "end_ele_handler"
    character_data_handler = "character_data_handler"
    processing_inst_handler = "processing_inst_handler"
    comment_handler = "comment_handler"
    start_cdata_section_handler = "start_cdata_section_handler"
    end_cdata_section_handler = "end_cdata_section_handler"
    default_handler = "default_handler"
)

type NullHandler struct {
    name string
}
var null_handler = NullHandler{}

type StartElementHandler func(interface{}, string, map[string]string)
type EndElementHandler func(interface{}, string)
type CharacterDataHandler func(interface{}, string)
type PIHandler func(interface{}, string, string)
type CommentHandler func(interface{}, string)
type StartCDataSectionHandler func(interface{})
type EndCDataSectionHandler func(interface{})
type DefaultHandler func(interface{}, string)
//type StartNSDeclHandler func(interface{}, string, string)
//type EndNSDeclHandler func(interface{}, string)


type XmlParserHooker struct {
    handlerMap map[string]interface{}
}
//export InternalEndElementHandler
func InternalEndElementHandler(userData unsafe.Pointer, name *C.XML_Char) {
    ud := (*userDataStructure)(userData)
    cname := (*C.char)(name)

    handler, ok := ud.hooker.handlerMap[end_ele_handler]
    if !ok || handler == nil {
        log.Print( "end element handler not defined" )
        return
    }

    finalHandler, ok := handler.(EndElementHandler)
    if !ok {
        log.Print( "end element handler type error" )
        return
    }

    finalHandler(ud.data, C.GoString(cname))
}

//export InternalStartElementHandler
func InternalStartElementHandler(userData unsafe.Pointer,  name *C.XML_Char, attrs **C.XML_Char) {
    ud := (*userDataStructure)(userData)
    cname := (*C.char)(name)
    cattrs := (**C.char)(unsafe.Pointer(attrs))

    arrLen := C.getExpatArrayLen(cattrs)

    var tmpSlice []*C.char
    sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&tmpSlice)))
    sliceHeader.Cap = int(arrLen)
    sliceHeader.Len = int(arrLen)
    sliceHeader.Data = uintptr(unsafe.Pointer(cattrs))

    attrsMap := make(map[string]string)
    for i := 0; i < len(tmpSlice); {
        key := C.GoString(tmpSlice[i])
        val := C.GoString(tmpSlice[i+1])
        attrsMap[key] = val
        i = i + 2
    }

    handler, ok := ud.hooker.handlerMap[start_ele_handler]
    if !ok || handler == nil {
        log.Println( "start element handler not defined" )
        return
    }

    finalHandler, ok := handler.(StartElementHandler)
    if !ok {
        log.Println( "start element handler type error" )
        return
    }

    finalHandler(ud.data, C.GoString(cname), attrsMap)
}

func convertCBytesArrayToGoString( bytesPtr *C.char, length C.int) string {
    var tmpSlice []byte
    sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&tmpSlice))
    sliceHeader.Cap = int(length)
    sliceHeader.Len = int(length)
    sliceHeader.Data = uintptr(unsafe.Pointer(bytesPtr))

    return string(tmpSlice[:int(length)])
}

//export InternalCharacterDataHandler
func InternalCharacterDataHandler(userData unsafe.Pointer, data *C.XML_Char, length C.int) {
    ud := (*userDataStructure)(userData)
    cdata := (*C.char)(data)
    goString := convertCBytesArrayToGoString(cdata, length)

    handler, ok := ud.hooker.handlerMap[character_data_handler]
    if !ok || handler == nil {
        log.Println("character data handler not defined")
        return
    }

    finalHandler, ok := handler.(CharacterDataHandler)
    if !ok {
        log.Println("character data handler type error")
        return;
    }

    finalHandler(ud.data, goString)
}

//export InternalProcessingInstHandler
func InternalProcessingInstHandler(userData unsafe.Pointer, target *C.XML_Char, data *C.XML_Char) {
    ud := (*userDataStructure)(userData)
    ctarget := (*C.char)(target)
    cdata := (*C.char)(data)

    handler, ok := ud.hooker.handlerMap[processing_inst_handler]
    if !ok || handler == nil {
        log.Println("processing inst handler not defined")
        return
    }

    finalHandler, ok := handler.(PIHandler)
    if !ok {
        log.Println("processing inst handler type error")
        return
    }

    finalHandler(ud.data, C.GoString(ctarget), C.GoString(cdata))
}

//export InternalStartCDataSectionHandler
func InternalStartCDataSectionHandler(userData unsafe.Pointer) {
    ud := (*userDataStructure)(userData)
    handler, ok := ud.hooker.handlerMap[start_cdata_section_handler]
    if !ok || handler == nil {
        log.Println("start cdata section handler not defined")
        return
    }

    finalHandler, ok := handler.(StartCDataSectionHandler)
    if !ok {
        log.Println("start cdata section handler type error")
        return
    }

    finalHandler(ud.data)
}

//export InternalEndCDataSectionHandler
func InternalEndCDataSectionHandler(userData unsafe.Pointer) {
    ud := (*userDataStructure)(userData)
    handler, ok := ud.hooker.handlerMap[end_cdata_section_handler]
    if !ok || handler == nil {
        log.Println("end cdata section handler not defined")
        return
    }

    finalHandler, ok := handler.(EndCDataSectionHandler)
    if !ok {
        log.Println("end cdata section handler type error")
        return
    }

    finalHandler(ud.data)
}

//export InternalCommentHandler
func InternalCommentHandler(userData unsafe.Pointer, data *C.XML_Char) {
    ud := (*userDataStructure)(userData)
    cdata := (*C.char)(data)
    handler, ok := ud.hooker.handlerMap[comment_handler]
    if !ok || handler == nil {
        log.Println("comment handler not defined")
        return
    }

    finalHandler, ok := handler.(CommentHandler)
    if !ok {
        log.Println("comment handler type error")
        return
    }

    finalHandler(ud.data, C.GoString(cdata))
}

//export InternalDefaultHandler
func InternalDefaultHandler(userData unsafe.Pointer, data *C.XML_Char, length C.int) {
    ud := (*userDataStructure)(userData)
    cdata := (*C.char)(data)
    goString := convertCBytesArrayToGoString(cdata, length)

    handler, ok := ud.hooker.handlerMap[default_handler]
    if !ok || handler == nil {
        log.Println("default handler not defined")
        return
    }

    finalHandler, ok := handler.(DefaultHandler)
    if !ok {
        log.Println("default handler type error")
        return
    }

    finalHandler(ud.data, goString)
}

//Hook hook and unhook the callback specified by handler. If handler is nil, unhook will
//be performed. Otherwise, hook is performed.
func (self *XmlParserHooker) Hook(parser *XmlParser, handler interface{}) error {
    if parser == nil {
        return errors.New( "parameters cannot be nil" )
    }

    nullHandler, ok := handler.(NullHandler)
    if ok {
        unhookHandler, ok := self.handlerMap[nullHandler.name]
        if !ok {
            return errors.New( "no registered handler" )
        }
        var key string
        switch unhookHandler.(type) {
        case StartElementHandler:
            C.unhookStartElementHandler(parser.parserHandler)
            key = start_ele_handler
        case EndElementHandler:
            C.unhookEndElementHandler(parser.parserHandler)
            key = end_ele_handler
        case CharacterDataHandler:
            C.unhookCharacterDataHandler(parser.parserHandler)
            key = character_data_handler
        case PIHandler:
            C.unhookPIHandler(parser.parserHandler)
            key = processing_inst_handler
        case CommentHandler:
            C.unhookCommentHandler(parser.parserHandler)
            key = comment_handler
        case StartCDataSectionHandler:
            C.unhookStartCDataSectionHandler(parser.parserHandler)
            key = start_cdata_section_handler
        case EndCDataSectionHandler:
            C.unhookEndCDataSectionHandler(parser.parserHandler)
            key = end_cdata_section_handler
        case DefaultHandler:
            C.unhookDefaultHandler(parser.parserHandler)
            key = default_handler
        default:
            return errors.New( "unsupported handler type" )
        }
        delete(self.handlerMap, key)
    } else {
        var key string
        switch handler.(type) {
        case StartElementHandler:
            C.hookStartElementHandler(parser.parserHandler)
            key = start_ele_handler
        case EndElementHandler:
            C.hookEndElementHandler(parser.parserHandler)
            key = end_ele_handler
        case CharacterDataHandler:
            C.hookCharacterDataHandler(parser.parserHandler)
            key = character_data_handler
        case PIHandler:
            C.hookPIHandler(parser.parserHandler)
            key = processing_inst_handler
        case CommentHandler:
            C.hookCommentHandler(parser.parserHandler)
            key = comment_handler
        case StartCDataSectionHandler:
            C.hookStartCDataSectionHandler(parser.parserHandler)
            key = start_cdata_section_handler
        case EndCDataSectionHandler:
            C.hookEndCDataSectionHandler(parser.parserHandler)
            key = end_cdata_section_handler
        case DefaultHandler:
            C.hookDefaultHandler(parser.parserHandler)
            key = default_handler
        default:
            return errors.New( "unsupported handler type" )
        }
        self.handlerMap[key] = handler
    }

    return nil
}
