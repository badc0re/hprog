// HCore.g4
grammar HCore;

// Lexer
OP: '(';
CP: ')';

// Operators
MUL: '*';
DIV: '/';
ADD: '+';
SUB: '-';
NUMBER: [0-9]+;
WHITESPACE: [ \r\n\t]+ -> skip;

// Parser
start : expression* EOF;


operator :
   (MUL | DIV | ADD | SUB)
   ;

expression 
   : 
   NUMBER | (OP operator expression+ CP)
   ;
   
