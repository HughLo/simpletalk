package xmlparser

/*
#cgo LDFLAGS: -lexpat
#include <expat.h>
extern int getExpatArrayLen( char** data );
extern void hookStartElementHandler(XML_Parser parser);
extern void unhookStartElementHandler(XML_Parser parser);
extern void hookEndElementHandler(XML_Parser parser);
extern void unhookEndElementHandler(XML_Parser parser);
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
)

type NullHandler struct {
    name string
}
var null_handler = NullHandler{}

type StartElementHandler func(interface{}, string, map[string]string)
type EndElementHandler func(interface{}, string)

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
        default:
            return errors.New( "unsupported handler type" )
        }
        self.handlerMap[key] = handler
    }

    return nil
}
