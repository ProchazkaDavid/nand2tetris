// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.


    // R0 = which color to draw (0 ~ white, -1 ~ black)
    @R0
    M=0

    // DRAW white if KBD == 0
    @KBD
    D=M
    @DRAW
    D;JEQ

    // DRAW black
    @R0
    M=-1
    
(DRAW) // params: R0 ~ color to be drawn
    // i = SCREEN + 8191
    @SCREEN
    D=A
    @8191
    D=D+A
    @i
    M=D

    // while i >= 0:
    //     RAM[i] = R0
    //     i--
    // goto CHOOSECOLOR
(LOOP)
    // if i < 0 goto CHOOSECOLOR
    @i
    D=M
    @CHOOSECOLOR
    D;JLT

    // RAM[i] = R0
    @R0
    D=M
    @i
    A=M
    M=D

    // i--
    @i
    M=M-1

    // goto LOOP
    @LOOP
    0;JMP

    // goto (WHITE if R0 == 0 else BLACK)
(CHOOSECOLOR)
    @R0
    D=M
    @WHITE
    D;JEQ
    @BLACK
    0;JMP

(WHITE)
    // goto (WHITE if KBD == 0 else DRAW black)
    @KBD
    D=M
    @WHITE
    D;JEQ

    // set black color to be drawn
    @R0
    M=-1

    // goto DRAW
    @DRAW
    0;JMP

(BLACK)
    // goto (BLACK if KBD != 0 else DRAW white)
    @KBD
    D=M
    @BLACK
    D;JNE

    // set white color to be drawn
    @R0
    M=0

    // goto DRAW
    @DRAW
    0;JMP
