# Elecraft CLI and library

Evolving suite of library components and cli for KX2, KX3, and K3s. Providing :
- Basic Serial connection  
- Basic command Send and Response receive
- Framework for adding more commands with parsing of responses.

The main goal of this project was to use the elecraft rigs as CW practice tools.
As such, functions are added around these capabilities first :
- CLI to send keyed CW to stdout (copying serial characters to stdout after executing a TT1; command)
- Terminal based UI to practice sending CW by using a text as a model (dictation) and verifying the keyed CW against that text.
- More to come...

## CW Dictation
The provided text is filtered to remove keeping the most used characters (alphabet, numeric, space, comma, period, question mark), and translated to upper case.

You can position the cursor where ever you want the dictation to start, using the space bar to move a page down, arrows to move up, down, and righ/left.
Because space is automatically inserted by the keyer, timing is critical, specially for comma and period.
The pro-sign BT advance the cursor by one if you need to - however, your keyboard arrows work too. This would be specially useful if you are using a justified text with multiple spaces which you cannot do with the keyer.

If using a KXPA100 or PX3/P3 in combination with the KX2/3 - use the buffered mode (-b).