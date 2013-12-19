#include <stdlib.h>
#include <stdio.h>
#include "_cgo_export.h"

int getExpatArrayLen( char** data )
{
    int len = 0;
    while( data[len] != NULL )
    {
        ++len;
    }
    return len;
}

void startElementHandler( void *userData, const XML_Char *name, const XML_Char **attrs)
{
    InternalStartElementHandler(userData, (XML_Char*)name, (XML_Char**)attrs); 
}

void endElementHandler( void *userData, const XML_Char *name )
{
    InternalEndElementHandler(userData, (XML_Char*)name);
}

void characterDataHandler(void *userData, const XML_Char *s, int len)
{
    InternalCharacterDataHandler(userData, (XML_Char*)s, len);    
}

void processingInstHandler(void *userData, const XML_Char *target, const XML_Char *data)
{
    InternalProcessingInstHandler(userData, (XML_Char*)target, (XML_Char*)data);
}

void commentHandler(void *userData, const XML_Char *data)
{
    InternalCommentHandler(userData, (XML_Char*)data);
}

void startCDataSectionHandler(void *userData)
{
    InternalStartCDataSectionHandler(userData);
}

void endCDataSectionHandler(void *userData)
{
    InternalEndCDataSectionHandler(userData);
}

void defaultHandler(void *userData, const XML_Char *s, int len)
{
    InternalDefaultHandler(userData, (XML_Char*)s, len);
}

void hookStartElementHandler( XML_Parser parser )
{
    XML_SetStartElementHandler( parser, startElementHandler );
}

void unhookStartElementHandler( XML_Parser parser )
{
    XML_SetStartElementHandler(parser, NULL);
}

void hookEndElementHandler( XML_Parser parser )
{
    XML_SetEndElementHandler(parser, endElementHandler);
}

void unhookEndElementHandler( XML_Parser parser )
{
    XML_SetEndElementHandler(parser, NULL);
}

void hookCharacterDataHandler(XML_Parser parser)
{
    XML_SetCharacterDataHandler(parser, characterDataHandler);
}

void unhookCharacterDataHandler(XML_Parser parser)
{
    XML_SetCharacterDataHandler(parser, NULL);
}

void hookPIHandler(XML_Parser parser)
{
    XML_SetProcessingInstructionHandler(parser, processingInstHandler);
}

void unhookPIHandler(XML_Parser parser)
{
    XML_SetProcessingInstructionHandler(parser, NULL);
}

void hookCommentHandler(XML_Parser parser)
{
    XML_SetCommentHandler(parser, commentHandler);
}

void unhookCommentHandler(XML_Parser parser)
{
    XML_SetCommentHandler(parser, NULL);
}

void hookStartCDataSectionHandler(XML_Parser parser)
{
    XML_SetStartCdataSectionHandler(parser, startCDataSectionHandler);
}

void unhookStartCDataSectionHandler(XML_Parser parser)
{
    XML_SetStartCdataSectionHandler(parser, NULL);
}

void hookEndCDataSectionHandler(XML_Parser parser)
{
    XML_SetEndCdataSectionHandler(parser, endCDataSectionHandler);
}

void unhookEndCDataSectionHandler(XML_Parser parser)
{
    XML_SetEndCdataSectionHandler(parser, NULL);
}

void hookDefaultHandler(XML_Parser parser)
{
    XML_SetDefaultHandler(parser, defaultHandler);
}

void unhookDefaultHandler(XML_Parser parser)
{
    XML_SetDefaultHandler(parser, NULL);
}
