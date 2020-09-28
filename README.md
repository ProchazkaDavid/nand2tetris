# [nand2tetris](https://www.nand2tetris.org)

> Build a Modern Computer from First Principles.

## Table of Contents

1. [Compiler](#compiler)
2. [VM Translator](#vm-translator)
3. [Assembler](#assembler)
4. [Computer](#computer)

---

## [Compiler](./compiler)

- uses LL(2) parser to process `.jack` files and generate stack based bytecode (`.vm`)

### Example of `.jack` file

```
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/11/Average/Main.jack

// (Same as projects/09/Average/Main.jack)

// Inputs some numbers and computes their average
class Main {
   function void main() {
     var Array a; 
     var int length;
     var int i, sum;

     let length = Keyboard.readInt("How many numbers? ");
     let a = Array.new(length); // constructs the array
     
     let i = 0;
     while (i < length) {
        let a[i] = Keyboard.readInt("Enter a number: ");
        let sum = sum + a[i];
        let i = i + 1;
     }
     
     do Output.printString("The average is ");
     do Output.printInt(sum / length);
     return;
   }
}
```

[More examples of .hack files](./compiler/examples)

## [VM Translator](./vm)

- processes `.vm` files and generates assembly (`.asm`)

### Example of `.vm` file

```
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/08/ProgramFlow/BasicLoop/BasicLoop.vm

// Computes the sum 1 + 2 + ... + argument[0] and pushes the 
// result onto the stack. Argument[0] is initialized by the test 
// script before this code starts running.
push constant 0    
pop local 0         // initializes sum = 0
label LOOP_START
push argument 0    
push local 0
add
pop local 0	        // sum = sum + counter
push argument 0
push constant 1
sub
pop argument 0      // counter--
push argument 0
if-goto LOOP_START  // If counter > 0, goto LOOP_START
push local 0
```

[More examples of .vm files](./vm/examples)

## [Assembler](./assembler)

- processes `.asm` file and generates machine code (`.hack`)
- two-pass assembler
- two types of intructions
  - A-instruction for setting the address register to a 15-bit value
  - C-instruction for basic arithmetic and logical operations

### Example of `.asm` file

```
// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)

    
    // R2 = 0
    @R2
    M=0

    // i = R0
    @R0
    D=M
    @i
    M=D
    
    // while i > 0:
    //     R2 += R1
    //     i--
    // goto END
(LOOP)

    // if i == 0: goto END
    @i
    D=M
    @END
    D;JEQ

    // R2 += R1
    @R1
    D=M
    @R2
    M=M+D

    // i--
    @i
    M=M-1
    
    // goto LOOP
    @LOOP
    0;JMP

(END)
    @END
    0;JMP
```

[More examples of .asm files](./assembler/examples)

## [Computer](./computer)

- 16-bit computer
- Harvard architecture 
- 3 CPU registers (Address, Data, Memory)
- 32K ROM
- 16K RAM with two memory-mapped I/O devices: a screen (8K) and a keyboard (1 word).

### Memory Chip Example

```
/**
 * The complete address space of the Hack computer's memory,
 * including RAM and memory-mapped I/O. 
 * The chip facilitates read and write operations, as follows:
 *     Read:  out(t) = Memory[address(t)](t)
 *     Write: if load(t-1) then Memory[address(t-1)](t) = in(t-1)
 * In words: the chip always outputs the value stored at the memory 
 * location specified by address. If load==1, the in value is loaded 
 * into the memory location specified by address. This value becomes 
 * available through the out output from the next time step onward.
 * Address space rules:
 * Only the upper 16K+8K+1 words of the Memory chip are used. 
 * Access to address>0x6000 is invalid. Access to any address in 
 * the range 0x4000-0x5FFF results in accessing the screen memory 
 * map. Access to address 0x6000 results in accessing the keyboard 
 * memory map. The behavior in these addresses is described in the 
 * Screen and Keyboard chip specifications given in the book.
 */

CHIP Memory {
    IN in[16], load, address[15];
    OUT out[16];

    PARTS:
    DMux(in=load, sel=address[14], a=loadRam, b=loadScreen);

    RAM16K(in=in, load=loadRam, address=address[0..13], out=ramResult);
    Screen(in=in, load=loadScreen, address=address[0..12], out=screenResult);
    Keyboard(out=keyboardResult);

    Mux4Way16(a=ramResult, b=ramResult, c=screenResult, d=keyboardResult, sel=address[13..14], out=out);
}
```

[Computer parts](./computer)
