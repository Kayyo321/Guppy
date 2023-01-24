# Guppy

<p align="center">A smart, and speedy language that compiles to machine code.</p>
<p align="center">
<img src="logo.png" alt="Guppy Lang Logo" width="480" height="480">
</p>

# Hello World; From Guppy

<p align="center">
Guppy is a high-level language that compiles down to machine code. Guppy has no 
runtime, and doesn't use a garbage collector. Functions like `malloc`, `calloc`,
and `free` to deal with heap objects. Guppy has bindings for the C Standard Library,
meaning it needs the GNU Compiler Collection to be installed on the system when compiling
with Guppy; it also requires NASM (Netwide Assembler) to convert the generated assembly 
(based on the source code) to `.obj` (object) files.

To download these tools, refer to the following websites:

* <a href="https://www.nasm.us/">NASM (Netwide Assembler).us</a>
* <a href="https://www.mingw-w64.org/">MinGW-w64 (GCC).org</a>

# Installation

## Option 1

<div align="center">
You can download the Guppy source code and build it using the
<a href="https://go.dev">Go Compiler</a>.

:warning: Documentation may not be there if you happen to download between an update, or bux fixes. :warning:

:warning: You will also have to add this to your Environment Variables :warning:
</div>  

## Option 2
<div align="center">
You can download the latest working build from the <a href="https://github.com/Kayyo321/Guppy/releases">Releases</a> section.

:warning: You will still have to build this for your system, make sure you have the <a href="https://go.dev">Go Compiler</a> installed. :warning:

:warning: You will also have to add this to your Environment Variables :warning:
</div>

## Option 3
<div align="center">
You can install the <a href="https://google.com">Guppy Installer</a>, and Guppy will be installed correctly for your system.
</div>

# Use

<div align="center">
To learn the Guppy language go to the <a href="https://google.com">Docs</a> section.

If you want to test that your installation is working, try the following test:
</div>

* Copy the following code into a source file:

`
package main

import io

io | printfln("Hello, World!")
`

* Run the following commands in the cwd:

`
$ gup build -o myProg myCode.gpy
$ ./myProg
`

If when you run the file you get a simple "Hello, World!" program, then your build of Guppy is installed correctly!
If not, try reinstalling your build, or using another means of installation (i.e. building, or using the installer).

# Libraries

Guppy installs with its complete GSL (Guppy Standard Library) by default.
These include (but are not limited to)

| Library | Description                                                   |
|---------|---------------------------------------------------------------|
| io      | Input and Output for the system                               |
| std     | Holds standard functions && error handling                    |
 | convert | Allows conversion for standard types                          | 
 | math    | Contains many arithmetic functions and values                 |
 | time    | Bindings for dealing with delays, and getting date values     |
 | windows | Bindings for C's `windows.h` for creating && handling windows |
