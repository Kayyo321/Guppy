# Guppy

<p align="center">A smart, and speedy language that compiles to machine code.</p>

# Installation

## Option 1

<div align="center">
You can download the Guppy source code and build it using the
<a href="https://go.dev">Go Compiler</a>.

:warning: Documentation may not be there if you happen to download between an update, or bux fixes. :warning:
</div>  

## Option 2
<div align="center">
You can download the latest working build from the <a href"https://github.com/Kayyo321/Guppy/releases">Releases</a> section.

:warning: You will still have to build this for your system, make sure you have the <a href="https://go.dev">Go Compiler</a> installed.
</div>

## Option 3
<div align="center">
You can install the <a href="https://google.com">Guppy Installer</a>, and Guppy will be installed correctly for your system.
</div>

# Use

<div align="center">
To learn the Guppy language go to the <a href="https://google.com">Docs</a> section.

If you want to test that your install is working, try the following test:
</div>

1. Copy the following code into a source file:

```
package main

import io

io | printfln("Hello, World!")
```

2. Run the following commands in the cwd:

```
$ gup build -o myProg myCode.gpy
$ ./myProg
```

If when you run the file you get a simple "Hello, World!" program, then your build of Guppy is installed correctly!
If not, try reinstalling your build, or using another means of installation (i.e. building, or using the installer).
