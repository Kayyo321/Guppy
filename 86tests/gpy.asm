; nasm -fwin64 gpy.asm
; gcc -m64 -mconsole gpy.obj -o gpy.exe
; gpy

global  main
extern  puts

section .data
	message db 'Hello, World!', 10, 0

section .text

main:
    push rbp
    mov rbp, rsp
    sub rsp, 8*4   

    mov rcx, message
    call puts
   
    mov rax,0
    add rsp,8*4
    pop rbp
    ret