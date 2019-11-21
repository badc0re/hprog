# hprog

Why not as Lisp?
- I simply can't follow all the brackets

Not working
Lexer:
- 1..10 (ranges)
- .11 floats starting with dot
- 11+ numbers with stuff after them

Sample errors:
- no identifier, no sign, only number
hprog> x<2
{22 2 2 0} <NUMBER>
{34  4 0} \0

hprog> x < 2
{17  2 0} <
{22 2 4 0} <NUMBER>
{34  6 0} \0
