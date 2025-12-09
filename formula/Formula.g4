// Formula.g4
grammar Formula;

// Tokens
OP: '(' ;
CP: ')' ;
AND: '&' ;
OR: '|' ;
IMPLIES: '->' ;
BICONDITIONAL: '<->' ;
NOR:  '!|' ;
NAND:  '!&' ;
XOR:  '^' ;
NOT: '!' ;
VARIABLE: [a-zA-Z_0-9]+ ;
WHITESPACE: [ \t\r\n]+ -> skip ;

// Rules
start : expression EOF;

expression
    : OP left=expression op=(AND | OR | IMPLIES | BICONDITIONAL | NOR | NAND | XOR) right=expression CP #Binary
    | NOT negated=expression                              #Negation
    | VARIABLE                                            #Letter
    ;