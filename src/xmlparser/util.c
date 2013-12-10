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

void startElementHandler( void* userData, const XML_Char* name, const XML_Char** attrs)
{
    InternalStartElementHandler(userData, (XML_Char*)name, (XML_Char**)attrs); 
}

void endElementHandler( void* userData, const XML_Char* name )
{
    InternalEndElementHandler(userData, (XML_Char*)name);
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
