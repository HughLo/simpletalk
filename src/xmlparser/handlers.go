package xmlparser

/*
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

type Hooker interface {
    Hook(parser interface{}, handler interface{}) error
}

type StartElementHandler interface {
    HandleStartEle(userData interface{}, name string, attrs map[string]string)
}

type EndElementHandler interface {
    HandleEndEle(userData interface{}, name string)
}


type XmlParserHooker struct {
}

//export InternalEndElementHandler
func InternalEndElementHandler(userData unsafe.Pointer, name *C.XML_Char) {
    parser := (*XmlParser)(userData)
    cname := (*C.char)(name)

    handler, ok := parser.handlerMap["end_ele_handler"]
    if !ok || handler == nil {
        log.Print( "end element handler not defined" )
        return
    }

    finalHandler, ok := handler.(EndElementHandler)
    if !ok {
        log.Print( "end element handler type error" )
        return
    }

    finalHandler.HandleEndEle(parser, C.GoString(cname))
}

//export InternalStartElementHandler
func InternalStartElementHandler(userData unsafe.Pointer,  name *C.XML_Char, attrs **C.XML_Char) {
    parser := (*XmlParser)(userData)
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

    handler, ok := parser.handlerMap["start_ele_handler"]
    if !ok || handler == nil {
        log.Println( "start element handler not defined" )
        return
    }

    finalHandler, ok := handler.(StartElementHandler)
    if !ok {
        log.Println( "start element handler type error" )
        return
    }

    finalHandler.HandleStartEle(parser, C.GoString(cname), attrsMap)
}

func (self XmlParserHooker) Hook(parser interface{}, handler interface{}) error {
    if parser == nil {
        return errors.New( "parameters cannot be nil" )
    }

    myParser, ok := parser.(*XmlParser)
    if !ok {
        return errors.New( "the first argument shall be type of *XmlParser" )
    }

    switch handler.(type) {
    case StartElementHandler:
        if handler == nil {
            C.hookStartElementHandler(myParser.parserHandler)
        } else {
            C.unhookStartElementHandler(myParser.parserHandler)
        }
    case EndElementHandler:
        if handler == nil {
            C.hookEndElementHandler(myParser.parserHandler)
        } else {
            C.unhookEndElementHandler(myParser.parserHandler)
        }
    default:
        return errors.New( "unsupported handler type" )
    }

    return nil
}
