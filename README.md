# YODK - YOLOL Development Kit
[![Build Status](https://travis-ci.org/dbaumgarten/yodk.svg?branch=master)](https://travis-ci.org/dbaumgarten/yodk)

# What is YOLOL?
[YOLOL](https://wiki.starbasegame.com/index.php/YOLOL) is the ingame programming language for the upcoming game starbase.

# What is the YODK?
The yodk aims to be a toolkit that helps with the development of YOLOL-Scripts. YOLOL is a pretty limmited language and the lack of common features is really annoying for experienced programmers. The yodk (and in the future especially nolol) will try to mitigate this.

# Installation

## Binaries
You can find pre-build versions of the binaries [here](https://github.com/dbaumgarten/yodk/releases).
Just download them, unpack the zip file and place the binary somewhere in your PATH.

## From source
You will need to have the go-toolchain installed.  
```go install github.com/dbaumgarten/yodk```  
Binary releases will be provided later.  
Note: It is helpfull to add the yodk-binary to your path.

# Features

## Debugger
Run YOLOL-Code, set breakpoints, step through execution and inspect variables. Can run/debug multiple scripts in parallel.

```
# yodk debug fizbuzz.yolol 
Loaded and paused programm. Enter 'c' to execute
>>> l
--Programm--
>   1 :out="" a=0
    2 if not (a<=100) then goto 9 end
    3 if not (a%3==0 and a%5==0) then goto 5 end :out+="fizz"+"buzz"+" "
    4 goto 8
    5 if not (a%3==0) then goto 6 end :out+="fizz"+" " goto 8
    6 if not (a%5==0) then goto 7 end :out+="buzz"+" " goto 8
    7 :out+=a+" "
    8 a++ goto 2
    9 
>>> r
--Started--
>>> --Program finished--
--Enter r to restart--
>>> 
>>> vars
--Variables--
:out 'fizzbuzz 1 2 fizz 4 buzz fizz 7 8 fizz buzz 11 fizz 13 14 fizzbuzz 16 17 fizz 19 buzz fizz 22 23 fizz buzz 26 fizz 28 29 
>>>  

```

## Optimizer
Automatically optimize YOLOL-Code. Evaluate static expressions, shorten variable names and remove comments.

input.yolol:
```
superVerboseVariable = 5 //helpfull comment
anotherVar = (2+3)^2 * 5
:output = superVerboseVariable+anotherVar
```
```
# yodk optimize input.yolol
```
input-opt.yolol
```
a=5
b=125
:output=a+b
```

## Verifier
Check if YOLOL-code is valid.

input.yolol:
```
x = 3
if x == 4 then y=5
```
```
# yodk verify unopt-opt.yolol 
Error when parsing file: Parser error at Line: 2, Coloumn: 1 (up to Line: 2, Coloumn: 22): Expected 'end' after if statement. Found Token: ''(Newline)
```

## Formatter
Automatically format your yolol and nolol-code via cli or directly in your IDE.

## Testing of YOLOL-Code
Run custom testcases for your yolol scripts. Define testcases as yaml files (as seen below) and run them against your scripts.  
testcase.yaml:
```
scripts: 
  # list of scripts to run
  - name: fizbuzz.yolol
    # optional. Number of iterations for this script. Default: 1
    iterations: 1
    # optional. Number of lines to execute maximally for this script. Default: 0 = infinite
    maxlines: 10000
cases:
  # list of testcases
  - name: TestOutput
    # global variables to set before running. ':' can be omitted
    inputs:
      number: 0
    # expected value for global variables after running
    outputs:
      out: "fizzbuzz fizz buzz fizz fizz buzz fizz fizzbuzz fizz buzz fizz fizz buzz fizz fizzbuzz fizz buzz fizz fizz buzz fizz fizzbuzz fizz buzz fizz fizz buzz fizz fizzbuzz fizz buzz fizz fizz buzz fizz fizzbuzz fizz buzz fizz fizz buzz fizz fizzbuzz fizz buzz fizz fizz buzz "
      number: 101
  - name: TestOutput2
    inputs:
      number: 99
    outputs:
       out: "fizz buzz "
```
```
# yodk test testcase.yaml 
Running file: testcase.yaml
  Running case: TestOutput
  Running case: TestOutput2
Tests OK

```

Tests can also be debugged to find out what exactly is going wrong. Debugging happens with the usual debug command.

```
# yodk debug testcase.yaml
```


## vscode Extension
Extension for visual studio code. Adds syntax highlighting, validation, formatting and commands for compiling and optimizing code. [More information](vscode-yolol/README.md)

![VSCODE-Extension](images/vscode-screenshot.png)

## Implementation of the language server protocol
Can be used to add support for YOLOL to many IDEs. Just run ```yodk langserv``` to start a server instance. Refer to the documentation of your IDE on how to integrate a language server into it.

## NOLOL
A highly experimental dialect of YOLOL. Nolol is for yolol what typescript is for javascript. It adds missing features like loops, labeled gotos, multiline ifs and compile-time constants and is compiled to plain YOLOL to be used within starbase. The compiled code is also optimized. Variable names are shortened and as many statements are merged into one line as possible, to get as much as possible of the 20 lines of a YOLOL-Chip.

fizzbuzz.nolol:
```
const fizz = "fizz"
const buzz = "buzz"
const sep = " "
const upto = 100

if :out == 0 then
	:out = ""
end
// main loop
while number<=upto do
	if number%3==0 and number%5==0 then
		:out += fizz+buzz+sep
		goto next //skip other cases
	end
	if number%3==0 then
		:out+=fizz+sep
	end
	if number%5==0 then
	:out+=buzz+sep
	end
	next> 
	number++ 
end
```
```
yodk compile fizzbuz.nolol
```

fizzbuzz.yolol:
```
if :out==0 then :out="" end
if not (a<=100) then goto 6 end
if a%3==0 and a%5==0 then :out+="fizzbuzz " goto 5 end
if a%3==0 then :out+="fizz " end if a%5==0 then :out+="buzz " end
a++ goto 2
```

# Upcoming features
- Implement debug adapter protocoll to enable debugging with vscode

# Compatibility guarantees
Absolutely none. There will be massive changes to the codebase in the near future and things WILL definetly break.  

Also, as starbase has not been released, there is a lot of guesswork involved in the code. The goal is to be 100% compatible to starbase's implementation.

# Contributing
Found bugs, suggestions, want to add features? Sure. Just open an issue.