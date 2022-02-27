"".T.M1 STEXT nosplit size=1 args=0x18 locals=0x0 funcid=0x0
	0x0000 00000 (main.go:19)	TEXT	"".T.M1(SB), NOSPLIT|ABIInternal, $0-24
	0x0000 00000 (main.go:19)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:19)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:19)	FUNCDATA	$5, "".T.M1.arginfo1(SB)
	0x0000 00000 (main.go:19)	RET
	0x0000 c3                                               .
"".T.M2 STEXT nosplit size=1 args=0x18 locals=0x0 funcid=0x0
	0x0000 00000 (main.go:20)	TEXT	"".T.M2(SB), NOSPLIT|ABIInternal, $0-24
	0x0000 00000 (main.go:20)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:20)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:20)	FUNCDATA	$5, "".T.M2.arginfo1(SB)
	0x0000 00000 (main.go:20)	RET
	0x0000 c3                                               .
"".main STEXT size=311 args=0x0 locals=0x88 funcid=0x0
	0x0000 00000 (main.go:31)	TEXT	"".main(SB), ABIInternal, $136-0
	0x0000 00000 (main.go:31)	LEAQ	-8(SP), R12
	0x0005 00005 (main.go:31)	CMPQ	R12, 16(R14)
	0x0009 00009 (main.go:31)	PCDATA	$0, $-2
	0x0009 00009 (main.go:31)	JLS	301
	0x000f 00015 (main.go:31)	PCDATA	$0, $-1
	0x000f 00015 (main.go:31)	SUBQ	$136, SP
	0x0016 00022 (main.go:31)	MOVQ	BP, 128(SP)
	0x001e 00030 (main.go:31)	LEAQ	128(SP), BP
	0x0026 00038 (main.go:31)	FUNCDATA	$0, gclocals·7d2d5fca80364273fb07d5820a76fef4(SB)
	0x0026 00038 (main.go:31)	FUNCDATA	$1, gclocals·5281ddd282d2ee46e64efba251339d2f(SB)
	0x0026 00038 (main.go:31)	FUNCDATA	$2, "".main.stkobj(SB)
	0x0026 00038 (main.go:37)	MOVQ	$17, ""..autotmp_15+104(SP)
	0x002f 00047 (main.go:37)	LEAQ	go.string."hello, interface"(SB), CX
	0x0036 00054 (main.go:37)	MOVQ	CX, ""..autotmp_15+112(SP)
	0x003b 00059 (main.go:37)	MOVQ	$16, ""..autotmp_15+120(SP)
	0x0044 00068 (main.go:37)	LEAQ	type."".T(SB), AX
	0x004b 00075 (main.go:37)	LEAQ	""..autotmp_15+104(SP), BX
	0x0050 00080 (main.go:37)	PCDATA	$1, $0
	0x0050 00080 (main.go:37)	CALL	runtime.convT2E(SB)
	0x0055 00085 (main.go:41)	MOVQ	AX, ""..autotmp_42+64(SP)
	0x005a 00090 (main.go:41)	MOVQ	BX, "".ei.data+48(SP)
	0x005f 00095 (main.go:40)	MOVQ	$17, ""..autotmp_15+104(SP)
	0x0068 00104 (main.go:40)	LEAQ	go.string."hello, interface"(SB), CX
	0x006f 00111 (main.go:40)	MOVQ	CX, ""..autotmp_15+112(SP)
	0x0074 00116 (main.go:40)	MOVQ	$16, ""..autotmp_15+120(SP)
	0x007d 00125 (main.go:40)	LEAQ	go.itab."".T,"".NonEmptyInterface(SB), AX
	0x0084 00132 (main.go:40)	LEAQ	""..autotmp_15+104(SP), BX
	0x0089 00137 (main.go:40)	PCDATA	$1, $1
	0x0089 00137 (main.go:40)	CALL	runtime.convT2I(SB)
	0x008e 00142 (main.go:42)	MOVQ	AX, ""..autotmp_43+56(SP)
	0x0093 00147 (main.go:42)	MOVQ	BX, "".i.data+40(SP)
	0x0098 00152 (main.go:41)	MOVUPS	X15, ""..autotmp_19+88(SP)
	0x009e 00158 (main.go:41)	MOVQ	""..autotmp_42+64(SP), CX
	0x00a3 00163 (main.go:41)	MOVQ	CX, ""..autotmp_19+88(SP)
	0x00a8 00168 (main.go:41)	MOVQ	"".ei.data+48(SP), CX
	0x00ad 00173 (main.go:41)	MOVQ	CX, ""..autotmp_19+96(SP)
	0x00b2 00178 (<unknown line number>)	NOP
	0x00b2 00178 ($GOROOT/src/fmt/print.go:274)	MOVQ	os.Stdout(SB), CX
	0x00b9 00185 ($GOROOT/src/fmt/print.go:274)	MOVL	$1, DI
	0x00be 00190 ($GOROOT/src/fmt/print.go:274)	MOVQ	DI, SI
	0x00c1 00193 ($GOROOT/src/fmt/print.go:274)	LEAQ	go.itab.*os.File,io.Writer(SB), AX
	0x00c8 00200 ($GOROOT/src/fmt/print.go:274)	MOVQ	CX, BX
	0x00cb 00203 ($GOROOT/src/fmt/print.go:274)	LEAQ	""..autotmp_19+88(SP), CX
	0x00d0 00208 ($GOROOT/src/fmt/print.go:274)	PCDATA	$1, $2
	0x00d0 00208 ($GOROOT/src/fmt/print.go:274)	CALL	fmt.Fprintln(SB)
	0x00d5 00213 (main.go:42)	MOVQ	""..autotmp_43+56(SP), CX
	0x00da 00218 (main.go:42)	TESTQ	CX, CX
	0x00dd 00221 (main.go:42)	JEQ	229
	0x00df 00223 (main.go:42)	MOVQ	8(CX), AX
	0x00e3 00227 (main.go:42)	JMP	232
	0x00e5 00229 (main.go:42)	MOVQ	CX, AX
	0x00e8 00232 (main.go:42)	MOVUPS	X15, ""..autotmp_23+72(SP)
	0x00ee 00238 (main.go:42)	MOVQ	AX, ""..autotmp_23+72(SP)
	0x00f3 00243 (main.go:42)	MOVQ	"".i.data+40(SP), DX
	0x00f8 00248 (main.go:42)	MOVQ	DX, ""..autotmp_23+80(SP)
	0x00fd 00253 (<unknown line number>)	NOP
	0x00fd 00253 ($GOROOT/src/fmt/print.go:274)	MOVQ	os.Stdout(SB), BX
	0x0104 00260 ($GOROOT/src/fmt/print.go:274)	LEAQ	go.itab.*os.File,io.Writer(SB), AX
	0x010b 00267 ($GOROOT/src/fmt/print.go:274)	LEAQ	""..autotmp_23+72(SP), CX
	0x0110 00272 ($GOROOT/src/fmt/print.go:274)	MOVL	$1, DI
	0x0115 00277 ($GOROOT/src/fmt/print.go:274)	MOVQ	DI, SI
	0x0118 00280 ($GOROOT/src/fmt/print.go:274)	PCDATA	$1, $0
	0x0118 00280 ($GOROOT/src/fmt/print.go:274)	CALL	fmt.Fprintln(SB)
	0x011d 00285 (main.go:43)	MOVQ	128(SP), BP
	0x0125 00293 (main.go:43)	ADDQ	$136, SP
	0x012c 00300 (main.go:43)	RET
	0x012d 00301 (main.go:43)	NOP
	0x012d 00301 (main.go:31)	PCDATA	$1, $-1
	0x012d 00301 (main.go:31)	PCDATA	$0, $-2
	0x012d 00301 (main.go:31)	CALL	runtime.morestack_noctxt(SB)
	0x0132 00306 (main.go:31)	PCDATA	$0, $-1
	0x0132 00306 (main.go:31)	JMP	0
	0x0000 4c 8d 64 24 f8 4d 3b 66 10 0f 86 1e 01 00 00 48  L.d$.M;f.......H
	0x0010 81 ec 88 00 00 00 48 89 ac 24 80 00 00 00 48 8d  ......H..$....H.
	0x0020 ac 24 80 00 00 00 48 c7 44 24 68 11 00 00 00 48  .$....H.D$h....H
	0x0030 8d 0d 00 00 00 00 48 89 4c 24 70 48 c7 44 24 78  ......H.L$pH.D$x
	0x0040 10 00 00 00 48 8d 05 00 00 00 00 48 8d 5c 24 68  ....H......H.\$h
	0x0050 e8 00 00 00 00 48 89 44 24 40 48 89 5c 24 30 48  .....H.D$@H.\$0H
	0x0060 c7 44 24 68 11 00 00 00 48 8d 0d 00 00 00 00 48  .D$h....H......H
	0x0070 89 4c 24 70 48 c7 44 24 78 10 00 00 00 48 8d 05  .L$pH.D$x....H..
	0x0080 00 00 00 00 48 8d 5c 24 68 e8 00 00 00 00 48 89  ....H.\$h.....H.
	0x0090 44 24 38 48 89 5c 24 28 44 0f 11 7c 24 58 48 8b  D$8H.\$(D..|$XH.
	0x00a0 4c 24 40 48 89 4c 24 58 48 8b 4c 24 30 48 89 4c  L$@H.L$XH.L$0H.L
	0x00b0 24 60 48 8b 0d 00 00 00 00 bf 01 00 00 00 48 89  $`H...........H.
	0x00c0 fe 48 8d 05 00 00 00 00 48 89 cb 48 8d 4c 24 58  .H......H..H.L$X
	0x00d0 e8 00 00 00 00 48 8b 4c 24 38 48 85 c9 74 06 48  .....H.L$8H..t.H
	0x00e0 8b 41 08 eb 03 48 89 c8 44 0f 11 7c 24 48 48 89  .A...H..D..|$HH.
	0x00f0 44 24 48 48 8b 54 24 28 48 89 54 24 50 48 8b 1d  D$HH.T$(H.T$PH..
	0x0100 00 00 00 00 48 8d 05 00 00 00 00 48 8d 4c 24 48  ....H......H.L$H
	0x0110 bf 01 00 00 00 48 89 fe e8 00 00 00 00 48 8b ac  .....H.......H..
	0x0120 24 80 00 00 00 48 81 c4 88 00 00 00 c3 e8 00 00  $....H..........
	0x0130 00 00 e9 c9 fe ff ff                             .......
	rel 3+0 t=24 type."".T+0
	rel 3+0 t=24 type."".T+0
	rel 3+0 t=24 type.*os.File+0
	rel 3+0 t=24 type.*os.File+0
	rel 50+4 t=15 go.string."hello, interface"+0
	rel 71+4 t=15 type."".T+0
	rel 81+4 t=7 runtime.convT2E+0
	rel 107+4 t=15 go.string."hello, interface"+0
	rel 128+4 t=15 go.itab."".T,"".NonEmptyInterface+0
	rel 138+4 t=7 runtime.convT2I+0
	rel 181+4 t=15 os.Stdout+0
	rel 196+4 t=15 go.itab.*os.File,io.Writer+0
	rel 209+4 t=7 fmt.Fprintln+0
	rel 256+4 t=15 os.Stdout+0
	rel 263+4 t=15 go.itab.*os.File,io.Writer+0
	rel 281+4 t=7 fmt.Fprintln+0
	rel 302+4 t=7 runtime.morestack_noctxt+0
"".NonEmptyInterface.M1 STEXT dupok size=112 args=0x10 locals=0x10 funcid=0x16
	0x0000 00000 (<autogenerated>:1)	TEXT	"".NonEmptyInterface.M1(SB), DUPOK|WRAPPER|ABIInternal, $16-16
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	68
	0x0006 00006 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0006 00006 (<autogenerated>:1)	SUBQ	$16, SP
	0x000a 00010 (<autogenerated>:1)	MOVQ	BP, 8(SP)
	0x000f 00015 (<autogenerated>:1)	LEAQ	8(SP), BP
	0x0014 00020 (<autogenerated>:1)	MOVQ	32(R14), R12
	0x0018 00024 (<autogenerated>:1)	TESTQ	R12, R12
	0x001b 00027 (<autogenerated>:1)	JNE	95
	0x001d 00029 (<autogenerated>:1)	NOP
	0x001d 00029 (<autogenerated>:1)	MOVQ	AX, ""..this+24(FP)
	0x0022 00034 (<autogenerated>:1)	MOVQ	BX, ""..this+32(FP)
	0x0027 00039 (<autogenerated>:1)	FUNCDATA	$0, gclocals·09cf9819fc716118c209c2d2155a3632(SB)
	0x0027 00039 (<autogenerated>:1)	FUNCDATA	$1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
	0x0027 00039 (<autogenerated>:1)	FUNCDATA	$5, "".NonEmptyInterface.M1.arginfo1(SB)
	0x0027 00039 (<autogenerated>:1)	MOVQ	AX, ""..this+24(SP)
	0x002c 00044 (<autogenerated>:1)	MOVQ	BX, ""..this+32(SP)
	0x0031 00049 (<autogenerated>:1)	MOVQ	24(AX), CX
	0x0035 00053 (<autogenerated>:1)	MOVQ	BX, AX
	0x0038 00056 (<autogenerated>:1)	PCDATA	$1, $1
	0x0038 00056 (<autogenerated>:1)	CALL	CX
	0x003a 00058 (<autogenerated>:1)	MOVQ	8(SP), BP
	0x003f 00063 (<autogenerated>:1)	ADDQ	$16, SP
	0x0043 00067 (<autogenerated>:1)	RET
	0x0044 00068 (<autogenerated>:1)	NOP
	0x0044 00068 (<autogenerated>:1)	PCDATA	$1, $-1
	0x0044 00068 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0044 00068 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x0049 00073 (<autogenerated>:1)	MOVQ	BX, 16(SP)
	0x004e 00078 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x0053 00083 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x0058 00088 (<autogenerated>:1)	MOVQ	16(SP), BX
	0x005d 00093 (<autogenerated>:1)	PCDATA	$0, $-1
	0x005d 00093 (<autogenerated>:1)	JMP	0
	0x005f 00095 (<autogenerated>:1)	LEAQ	24(SP), R13
	0x0064 00100 (<autogenerated>:1)	CMPQ	(R12), R13
	0x0068 00104 (<autogenerated>:1)	JNE	29
	0x006a 00106 (<autogenerated>:1)	MOVQ	SP, (R12)
	0x006e 00110 (<autogenerated>:1)	JMP	29
	0x0000 49 3b 66 10 76 3e 48 83 ec 10 48 89 6c 24 08 48  I;f.v>H...H.l$.H
	0x0010 8d 6c 24 08 4d 8b 66 20 4d 85 e4 75 42 48 89 44  .l$.M.f M..uBH.D
	0x0020 24 18 48 89 5c 24 20 48 89 44 24 18 48 89 5c 24  $.H.\$ H.D$.H.\$
	0x0030 20 48 8b 48 18 48 89 d8 ff d1 48 8b 6c 24 08 48   H.H.H....H.l$.H
	0x0040 83 c4 10 c3 48 89 44 24 08 48 89 5c 24 10 e8 00  ....H.D$.H.\$...
	0x0050 00 00 00 48 8b 44 24 08 48 8b 5c 24 10 eb a1 4c  ...H.D$.H.\$...L
	0x0060 8d 6c 24 18 4d 39 2c 24 75 b3 49 89 24 24 eb ad  .l$.M9,$u.I.$$..
	rel 2+0 t=25 type."".NonEmptyInterface+96
	rel 56+0 t=10 +0
	rel 79+4 t=7 runtime.morestack_noctxt+0
"".NonEmptyInterface.M2 STEXT dupok size=112 args=0x10 locals=0x10 funcid=0x16
	0x0000 00000 (<autogenerated>:1)	TEXT	"".NonEmptyInterface.M2(SB), DUPOK|WRAPPER|ABIInternal, $16-16
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	68
	0x0006 00006 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0006 00006 (<autogenerated>:1)	SUBQ	$16, SP
	0x000a 00010 (<autogenerated>:1)	MOVQ	BP, 8(SP)
	0x000f 00015 (<autogenerated>:1)	LEAQ	8(SP), BP
	0x0014 00020 (<autogenerated>:1)	MOVQ	32(R14), R12
	0x0018 00024 (<autogenerated>:1)	TESTQ	R12, R12
	0x001b 00027 (<autogenerated>:1)	JNE	95
	0x001d 00029 (<autogenerated>:1)	NOP
	0x001d 00029 (<autogenerated>:1)	MOVQ	AX, ""..this+24(FP)
	0x0022 00034 (<autogenerated>:1)	MOVQ	BX, ""..this+32(FP)
	0x0027 00039 (<autogenerated>:1)	FUNCDATA	$0, gclocals·09cf9819fc716118c209c2d2155a3632(SB)
	0x0027 00039 (<autogenerated>:1)	FUNCDATA	$1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
	0x0027 00039 (<autogenerated>:1)	FUNCDATA	$5, "".NonEmptyInterface.M2.arginfo1(SB)
	0x0027 00039 (<autogenerated>:1)	MOVQ	AX, ""..this+24(SP)
	0x002c 00044 (<autogenerated>:1)	MOVQ	BX, ""..this+32(SP)
	0x0031 00049 (<autogenerated>:1)	MOVQ	32(AX), CX
	0x0035 00053 (<autogenerated>:1)	MOVQ	BX, AX
	0x0038 00056 (<autogenerated>:1)	PCDATA	$1, $1
	0x0038 00056 (<autogenerated>:1)	CALL	CX
	0x003a 00058 (<autogenerated>:1)	MOVQ	8(SP), BP
	0x003f 00063 (<autogenerated>:1)	ADDQ	$16, SP
	0x0043 00067 (<autogenerated>:1)	RET
	0x0044 00068 (<autogenerated>:1)	NOP
	0x0044 00068 (<autogenerated>:1)	PCDATA	$1, $-1
	0x0044 00068 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0044 00068 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x0049 00073 (<autogenerated>:1)	MOVQ	BX, 16(SP)
	0x004e 00078 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x0053 00083 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x0058 00088 (<autogenerated>:1)	MOVQ	16(SP), BX
	0x005d 00093 (<autogenerated>:1)	PCDATA	$0, $-1
	0x005d 00093 (<autogenerated>:1)	JMP	0
	0x005f 00095 (<autogenerated>:1)	LEAQ	24(SP), R13
	0x0064 00100 (<autogenerated>:1)	CMPQ	(R12), R13
	0x0068 00104 (<autogenerated>:1)	JNE	29
	0x006a 00106 (<autogenerated>:1)	MOVQ	SP, (R12)
	0x006e 00110 (<autogenerated>:1)	JMP	29
	0x0000 49 3b 66 10 76 3e 48 83 ec 10 48 89 6c 24 08 48  I;f.v>H...H.l$.H
	0x0010 8d 6c 24 08 4d 8b 66 20 4d 85 e4 75 42 48 89 44  .l$.M.f M..uBH.D
	0x0020 24 18 48 89 5c 24 20 48 89 44 24 18 48 89 5c 24  $.H.\$ H.D$.H.\$
	0x0030 20 48 8b 48 20 48 89 d8 ff d1 48 8b 6c 24 08 48   H.H H....H.l$.H
	0x0040 83 c4 10 c3 48 89 44 24 08 48 89 5c 24 10 e8 00  ....H.D$.H.\$...
	0x0050 00 00 00 48 8b 44 24 08 48 8b 5c 24 10 eb a1 4c  ...H.D$.H.\$...L
	0x0060 8d 6c 24 18 4d 39 2c 24 75 b3 49 89 24 24 eb ad  .l$.M9,$u.I.$$..
	rel 2+0 t=25 type."".NonEmptyInterface+104
	rel 56+0 t=10 +0
	rel 79+4 t=7 runtime.morestack_noctxt+0
type..eq."".T STEXT dupok size=95 args=0x10 locals=0x20 funcid=0x0
	0x0000 00000 (<autogenerated>:1)	TEXT	type..eq."".T(SB), DUPOK|ABIInternal, $32-16
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	68
	0x0006 00006 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0006 00006 (<autogenerated>:1)	SUBQ	$32, SP
	0x000a 00010 (<autogenerated>:1)	MOVQ	BP, 24(SP)
	0x000f 00015 (<autogenerated>:1)	LEAQ	24(SP), BP
	0x0014 00020 (<autogenerated>:1)	FUNCDATA	$0, gclocals·dc9b0298814590ca3ffc3a889546fc8b(SB)
	0x0014 00020 (<autogenerated>:1)	FUNCDATA	$1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
	0x0014 00020 (<autogenerated>:1)	FUNCDATA	$5, type..eq."".T.arginfo1(SB)
	0x0014 00020 (<autogenerated>:1)	MOVQ	(AX), DX
	0x0017 00023 (<autogenerated>:1)	CMPQ	(BX), DX
	0x001a 00026 (<autogenerated>:1)	JNE	56
	0x001c 00028 (<autogenerated>:1)	MOVQ	16(AX), CX
	0x0020 00032 (<autogenerated>:1)	MOVQ	8(BX), DX
	0x0024 00036 (<autogenerated>:1)	MOVQ	8(AX), AX
	0x0028 00040 (<autogenerated>:1)	CMPQ	16(BX), CX
	0x002c 00044 (<autogenerated>:1)	JNE	56
	0x002e 00046 (<autogenerated>:1)	MOVQ	DX, BX
	0x0031 00049 (<autogenerated>:1)	PCDATA	$1, $1
	0x0031 00049 (<autogenerated>:1)	CALL	runtime.memequal(SB)
	0x0036 00054 (<autogenerated>:1)	JMP	58
	0x0038 00056 (<autogenerated>:1)	XORL	AX, AX
	0x003a 00058 (<autogenerated>:1)	PCDATA	$1, $-1
	0x003a 00058 (<autogenerated>:1)	MOVQ	24(SP), BP
	0x003f 00063 (<autogenerated>:1)	ADDQ	$32, SP
	0x0043 00067 (<autogenerated>:1)	RET
	0x0044 00068 (<autogenerated>:1)	NOP
	0x0044 00068 (<autogenerated>:1)	PCDATA	$1, $-1
	0x0044 00068 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0044 00068 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x0049 00073 (<autogenerated>:1)	MOVQ	BX, 16(SP)
	0x004e 00078 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x0053 00083 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x0058 00088 (<autogenerated>:1)	MOVQ	16(SP), BX
	0x005d 00093 (<autogenerated>:1)	PCDATA	$0, $-1
	0x005d 00093 (<autogenerated>:1)	JMP	0
	0x0000 49 3b 66 10 76 3e 48 83 ec 20 48 89 6c 24 18 48  I;f.v>H.. H.l$.H
	0x0010 8d 6c 24 18 48 8b 10 48 39 13 75 1c 48 8b 48 10  .l$.H..H9.u.H.H.
	0x0020 48 8b 53 08 48 8b 40 08 48 39 4b 10 75 0a 48 89  H.S.H.@.H9K.u.H.
	0x0030 d3 e8 00 00 00 00 eb 02 31 c0 48 8b 6c 24 18 48  ........1.H.l$.H
	0x0040 83 c4 20 c3 48 89 44 24 08 48 89 5c 24 10 e8 00  .. .H.D$.H.\$...
	0x0050 00 00 00 48 8b 44 24 08 48 8b 5c 24 10 eb a1     ...H.D$.H.\$...
	rel 50+4 t=7 runtime.memequal+0
	rel 79+4 t=7 runtime.morestack_noctxt+0
"".(*T).M1 STEXT dupok nosplit size=63 args=0x8 locals=0x8 funcid=0x16
	0x0000 00000 (<autogenerated>:1)	TEXT	"".(*T).M1(SB), DUPOK|NOSPLIT|WRAPPER|ABIInternal, $8-8
	0x0000 00000 (<autogenerated>:1)	SUBQ	$8, SP
	0x0004 00004 (<autogenerated>:1)	MOVQ	BP, (SP)
	0x0008 00008 (<autogenerated>:1)	LEAQ	(SP), BP
	0x000c 00012 (<autogenerated>:1)	MOVQ	32(R14), R12
	0x0010 00016 (<autogenerated>:1)	TESTQ	R12, R12
	0x0013 00019 (<autogenerated>:1)	JNE	46
	0x0015 00021 (<autogenerated>:1)	NOP
	0x0015 00021 (<autogenerated>:1)	FUNCDATA	$0, gclocals·1a65e721a2ccc325b382662e7ffee780(SB)
	0x0015 00021 (<autogenerated>:1)	FUNCDATA	$1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
	0x0015 00021 (<autogenerated>:1)	FUNCDATA	$5, "".(*T).M1.arginfo1(SB)
	0x0015 00021 (<autogenerated>:1)	MOVQ	AX, ""..this+16(SP)
	0x001a 00026 (<autogenerated>:1)	TESTQ	AX, AX
	0x001d 00029 (<autogenerated>:1)	JEQ	40
	0x001f 00031 (<autogenerated>:1)	MOVQ	(SP), BP
	0x0023 00035 (<autogenerated>:1)	ADDQ	$8, SP
	0x0027 00039 (<autogenerated>:1)	RET
	0x0028 00040 (<autogenerated>:1)	PCDATA	$1, $1
	0x0028 00040 (<autogenerated>:1)	CALL	runtime.panicwrap(SB)
	0x002d 00045 (<autogenerated>:1)	XCHGL	AX, AX
	0x002e 00046 (<autogenerated>:1)	LEAQ	16(SP), R13
	0x0033 00051 (<autogenerated>:1)	CMPQ	(R12), R13
	0x0037 00055 (<autogenerated>:1)	JNE	21
	0x0039 00057 (<autogenerated>:1)	MOVQ	SP, (R12)
	0x003d 00061 (<autogenerated>:1)	JMP	21
	0x0000 48 83 ec 08 48 89 2c 24 48 8d 2c 24 4d 8b 66 20  H...H.,$H.,$M.f 
	0x0010 4d 85 e4 75 19 48 89 44 24 10 48 85 c0 74 09 48  M..u.H.D$.H..t.H
	0x0020 8b 2c 24 48 83 c4 08 c3 e8 00 00 00 00 90 4c 8d  .,$H..........L.
	0x0030 6c 24 10 4d 39 2c 24 75 dc 49 89 24 24 eb d6     l$.M9,$u.I.$$..
	rel 41+4 t=7 runtime.panicwrap+0
"".(*T).M2 STEXT dupok nosplit size=63 args=0x8 locals=0x8 funcid=0x16
	0x0000 00000 (<autogenerated>:1)	TEXT	"".(*T).M2(SB), DUPOK|NOSPLIT|WRAPPER|ABIInternal, $8-8
	0x0000 00000 (<autogenerated>:1)	SUBQ	$8, SP
	0x0004 00004 (<autogenerated>:1)	MOVQ	BP, (SP)
	0x0008 00008 (<autogenerated>:1)	LEAQ	(SP), BP
	0x000c 00012 (<autogenerated>:1)	MOVQ	32(R14), R12
	0x0010 00016 (<autogenerated>:1)	TESTQ	R12, R12
	0x0013 00019 (<autogenerated>:1)	JNE	46
	0x0015 00021 (<autogenerated>:1)	NOP
	0x0015 00021 (<autogenerated>:1)	FUNCDATA	$0, gclocals·1a65e721a2ccc325b382662e7ffee780(SB)
	0x0015 00021 (<autogenerated>:1)	FUNCDATA	$1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
	0x0015 00021 (<autogenerated>:1)	FUNCDATA	$5, "".(*T).M2.arginfo1(SB)
	0x0015 00021 (<autogenerated>:1)	MOVQ	AX, ""..this+16(SP)
	0x001a 00026 (<autogenerated>:1)	TESTQ	AX, AX
	0x001d 00029 (<autogenerated>:1)	JEQ	40
	0x001f 00031 (<autogenerated>:1)	MOVQ	(SP), BP
	0x0023 00035 (<autogenerated>:1)	ADDQ	$8, SP
	0x0027 00039 (<autogenerated>:1)	RET
	0x0028 00040 (<autogenerated>:1)	PCDATA	$1, $1
	0x0028 00040 (<autogenerated>:1)	CALL	runtime.panicwrap(SB)
	0x002d 00045 (<autogenerated>:1)	XCHGL	AX, AX
	0x002e 00046 (<autogenerated>:1)	LEAQ	16(SP), R13
	0x0033 00051 (<autogenerated>:1)	CMPQ	(R12), R13
	0x0037 00055 (<autogenerated>:1)	JNE	21
	0x0039 00057 (<autogenerated>:1)	MOVQ	SP, (R12)
	0x003d 00061 (<autogenerated>:1)	JMP	21
	0x0000 48 83 ec 08 48 89 2c 24 48 8d 2c 24 4d 8b 66 20  H...H.,$H.,$M.f 
	0x0010 4d 85 e4 75 19 48 89 44 24 10 48 85 c0 74 09 48  M..u.H.D$.H..t.H
	0x0020 8b 2c 24 48 83 c4 08 c3 e8 00 00 00 00 90 4c 8d  .,$H..........L.
	0x0030 6c 24 10 4d 39 2c 24 75 dc 49 89 24 24 eb d6     l$.M9,$u.I.$$..
	rel 41+4 t=7 runtime.panicwrap+0
os.(*File).close STEXT dupok size=86 args=0x8 locals=0x10 funcid=0x16
	0x0000 00000 (<autogenerated>:1)	TEXT	os.(*File).close(SB), DUPOK|WRAPPER|ABIInternal, $16-8
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	52
	0x0006 00006 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0006 00006 (<autogenerated>:1)	SUBQ	$16, SP
	0x000a 00010 (<autogenerated>:1)	MOVQ	BP, 8(SP)
	0x000f 00015 (<autogenerated>:1)	LEAQ	8(SP), BP
	0x0014 00020 (<autogenerated>:1)	MOVQ	32(R14), R12
	0x0018 00024 (<autogenerated>:1)	TESTQ	R12, R12
	0x001b 00027 (<autogenerated>:1)	JNE	69
	0x001d 00029 (<autogenerated>:1)	NOP
	0x001d 00029 (<autogenerated>:1)	FUNCDATA	$0, gclocals·1a65e721a2ccc325b382662e7ffee780(SB)
	0x001d 00029 (<autogenerated>:1)	FUNCDATA	$1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
	0x001d 00029 (<autogenerated>:1)	FUNCDATA	$5, os.(*File).close.arginfo1(SB)
	0x001d 00029 (<autogenerated>:1)	MOVQ	AX, ""..this+24(SP)
	0x0022 00034 (<autogenerated>:1)	MOVQ	(AX), AX
	0x0025 00037 (<autogenerated>:1)	PCDATA	$1, $1
	0x0025 00037 (<autogenerated>:1)	CALL	os.(*file).close(SB)
	0x002a 00042 (<autogenerated>:1)	MOVQ	8(SP), BP
	0x002f 00047 (<autogenerated>:1)	ADDQ	$16, SP
	0x0033 00051 (<autogenerated>:1)	RET
	0x0034 00052 (<autogenerated>:1)	NOP
	0x0034 00052 (<autogenerated>:1)	PCDATA	$1, $-1
	0x0034 00052 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0034 00052 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x0039 00057 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x003e 00062 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x0043 00067 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0043 00067 (<autogenerated>:1)	JMP	0
	0x0045 00069 (<autogenerated>:1)	LEAQ	24(SP), R13
	0x004a 00074 (<autogenerated>:1)	CMPQ	(R12), R13
	0x004e 00078 (<autogenerated>:1)	JNE	29
	0x0050 00080 (<autogenerated>:1)	MOVQ	SP, (R12)
	0x0054 00084 (<autogenerated>:1)	JMP	29
	0x0000 49 3b 66 10 76 2e 48 83 ec 10 48 89 6c 24 08 48  I;f.v.H...H.l$.H
	0x0010 8d 6c 24 08 4d 8b 66 20 4d 85 e4 75 28 48 89 44  .l$.M.f M..u(H.D
	0x0020 24 18 48 8b 00 e8 00 00 00 00 48 8b 6c 24 08 48  $.H.......H.l$.H
	0x0030 83 c4 10 c3 48 89 44 24 08 e8 00 00 00 00 48 8b  ....H.D$......H.
	0x0040 44 24 08 eb bb 4c 8d 6c 24 18 4d 39 2c 24 75 cd  D$...L.l$.M9,$u.
	0x0050 49 89 24 24 eb c7                                I.$$..
	rel 38+4 t=7 os.(*file).close+0
	rel 58+4 t=7 runtime.morestack_noctxt+0
go.cuinfo.packagename. SDWARFCUINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
""..inittask SNOPTRDATA size=32
	0x0000 00 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 24+8 t=1 fmt..inittask+0
go.info.fmt.Println$abstract SDWARFABSFCN dupok size=42
	0x0000 04 66 6d 74 2e 50 72 69 6e 74 6c 6e 00 01 01 11  .fmt.Println....
	0x0010 61 00 00 00 00 00 00 11 6e 00 01 00 00 00 00 11  a.......n.......
	0x0020 65 72 72 00 01 00 00 00 00 00                    err.......
	rel 0+0 t=23 type.[]interface {}+0
	rel 0+0 t=23 type.error+0
	rel 0+0 t=23 type.int+0
	rel 19+4 t=31 go.info.[]interface {}+0
	rel 27+4 t=31 go.info.int+0
	rel 37+4 t=31 go.info.error+0
go.string."hello, interface" SRODATA dupok size=16
	0x0000 68 65 6c 6c 6f 2c 20 69 6e 74 65 72 66 61 63 65  hello, interface
runtime.memequal64·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 runtime.memequal64+0
runtime.gcbits.01 SRODATA dupok size=1
	0x0000 01                                               .
type..namedata.*func()- SRODATA dupok size=9
	0x0000 00 07 2a 66 75 6e 63 28 29                       ..*func()
type.*func() SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 9b 90 75 1b 08 08 08 36 00 00 00 00 00 00 00 00  ..u....6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*func()-+0
	rel 48+8 t=1 type.func()+0
type.func() SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 f6 bc 82 f6 02 08 08 33 00 00 00 00 00 00 00 00  .......3........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00                                      ....
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*func()-+0
	rel 44+4 t=-32763 type.*func()+0
runtime.interequal·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 runtime.interequal+0
type..namedata.*main.NonEmptyInterface. SRODATA dupok size=25
	0x0000 01 17 2a 6d 61 69 6e 2e 4e 6f 6e 45 6d 70 74 79  ..*main.NonEmpty
	0x0010 49 6e 74 65 72 66 61 63 65                       Interface
type.*"".NonEmptyInterface SRODATA size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 f6 af 20 d1 08 08 08 36 00 00 00 00 00 00 00 00  .. ....6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*main.NonEmptyInterface.+0
	rel 48+8 t=1 type."".NonEmptyInterface+0
runtime.gcbits.02 SRODATA dupok size=1
	0x0000 02                                               .
type..namedata.M1. SRODATA dupok size=4
	0x0000 01 02 4d 31                                      ..M1
type..namedata.M2. SRODATA dupok size=4
	0x0000 01 02 4d 32                                      ..M2
type."".NonEmptyInterface SRODATA size=112
	0x0000 10 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0010 4a 37 15 f3 07 08 08 14 00 00 00 00 00 00 00 00  J7..............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 02 00 00 00 00 00 00 00 02 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 20 00 00 00 00 00 00 00  ........ .......
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 24+8 t=1 runtime.interequal·f+0
	rel 32+8 t=1 runtime.gcbits.02+0
	rel 40+4 t=5 type..namedata.*main.NonEmptyInterface.+0
	rel 44+4 t=5 type.*"".NonEmptyInterface+0
	rel 48+8 t=1 type..importpath."".+0
	rel 56+8 t=1 type."".NonEmptyInterface+96
	rel 80+4 t=5 type..importpath."".+0
	rel 96+4 t=5 type..namedata.M1.+0
	rel 100+4 t=5 type.func()+0
	rel 104+4 t=5 type..namedata.M2.+0
	rel 108+4 t=5 type.func()+0
type..eqfunc."".T SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 type..eq."".T+0
type..namedata.*main.T. SRODATA dupok size=9
	0x0000 01 07 2a 6d 61 69 6e 2e 54                       ..*main.T
type..namedata.*func(*main.T)- SRODATA dupok size=16
	0x0000 00 0e 2a 66 75 6e 63 28 2a 6d 61 69 6e 2e 54 29  ..*func(*main.T)
type.*func(*"".T) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 73 30 e0 de 08 08 08 36 00 00 00 00 00 00 00 00  s0.....6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*func(*main.T)-+0
	rel 48+8 t=1 type.func(*"".T)+0
type.func(*"".T) SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 5b c2 60 8a 02 08 08 33 00 00 00 00 00 00 00 00  [.`....3........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*func(*main.T)-+0
	rel 44+4 t=-32763 type.*func(*"".T)+0
	rel 56+8 t=1 type.*"".T+0
type.*"".T SRODATA size=104
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 8b 54 ca 19 09 08 08 36 00 00 00 00 00 00 00 00  .T.....6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 02 00 02 00  ................
	0x0040 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0060 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*main.T.+0
	rel 48+8 t=1 type."".T+0
	rel 56+4 t=5 type..importpath."".+0
	rel 72+4 t=5 type..namedata.M1.+0
	rel 76+4 t=26 type.func()+0
	rel 80+4 t=26 "".(*T).M1+0
	rel 84+4 t=26 "".(*T).M1+0
	rel 88+4 t=5 type..namedata.M2.+0
	rel 92+4 t=26 type.func()+0
	rel 96+4 t=26 "".(*T).M2+0
	rel 100+4 t=26 "".(*T).M2+0
type..namedata.*func(main.T)- SRODATA dupok size=15
	0x0000 00 0d 2a 66 75 6e 63 28 6d 61 69 6e 2e 54 29     ..*func(main.T)
type.*func("".T) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 8d 60 38 6a 08 08 08 36 00 00 00 00 00 00 00 00  .`8j...6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*func(main.T)-+0
	rel 48+8 t=1 type.func("".T)+0
type.func("".T) SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 63 a7 41 a9 02 08 08 33 00 00 00 00 00 00 00 00  c.A....3........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*func(main.T)-+0
	rel 44+4 t=-32763 type.*func("".T)+0
	rel 56+8 t=1 type."".T+0
type..namedata.n- SRODATA dupok size=3
	0x0000 00 01 6e                                         ..n
type..namedata.s- SRODATA dupok size=3
	0x0000 00 01 73                                         ..s
type."".T SRODATA size=176
	0x0000 18 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0010 b5 a4 ef 44 07 08 08 19 00 00 00 00 00 00 00 00  ...D............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 02 00 00 00 00 00 00 00 02 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 02 00 02 00 40 00 00 00 00 00 00 00  ........@.......
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0080 00 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0090 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x00a0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 24+8 t=1 type..eqfunc."".T+0
	rel 32+8 t=1 runtime.gcbits.02+0
	rel 40+4 t=5 type..namedata.*main.T.+0
	rel 44+4 t=5 type.*"".T+0
	rel 48+8 t=1 type..importpath."".+0
	rel 56+8 t=1 type."".T+96
	rel 80+4 t=5 type..importpath."".+0
	rel 96+8 t=1 type..namedata.n-+0
	rel 104+8 t=1 type.int+0
	rel 120+8 t=1 type..namedata.s-+0
	rel 128+8 t=1 type.string+0
	rel 144+4 t=5 type..namedata.M1.+0
	rel 148+4 t=26 type.func()+0
	rel 152+4 t=26 "".(*T).M1+0
	rel 156+4 t=26 "".T.M1+0
	rel 160+4 t=5 type..namedata.M2.+0
	rel 164+4 t=26 type.func()+0
	rel 168+4 t=26 "".(*T).M2+0
	rel 172+4 t=26 "".T.M2+0
runtime.nilinterequal·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 runtime.nilinterequal+0
type..namedata.*interface {}- SRODATA dupok size=15
	0x0000 00 0d 2a 69 6e 74 65 72 66 61 63 65 20 7b 7d     ..*interface {}
type.*interface {} SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 4f 0f 96 9d 08 08 08 36 00 00 00 00 00 00 00 00  O......6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*interface {}-+0
	rel 48+8 t=1 type.interface {}+0
type.interface {} SRODATA dupok size=80
	0x0000 10 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0010 e7 57 a0 18 02 08 08 14 00 00 00 00 00 00 00 00  .W..............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 24+8 t=1 runtime.nilinterequal·f+0
	rel 32+8 t=1 runtime.gcbits.02+0
	rel 40+4 t=5 type..namedata.*interface {}-+0
	rel 44+4 t=-32763 type.*interface {}+0
	rel 56+8 t=1 type.interface {}+80
type..namedata.*[]interface {}- SRODATA dupok size=17
	0x0000 00 0f 2a 5b 5d 69 6e 74 65 72 66 61 63 65 20 7b  ..*[]interface {
	0x0010 7d                                               }
type.*[]interface {} SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 f3 04 9a e7 08 08 08 36 00 00 00 00 00 00 00 00  .......6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*[]interface {}-+0
	rel 48+8 t=1 type.[]interface {}+0
type.[]interface {} SRODATA dupok size=56
	0x0000 18 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 70 93 ea 2f 02 08 08 17 00 00 00 00 00 00 00 00  p../............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*[]interface {}-+0
	rel 44+4 t=-32763 type.*[]interface {}+0
	rel 48+8 t=1 type.interface {}+0
go.itab."".T,"".NonEmptyInterface SRODATA dupok size=40
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 b5 a4 ef 44 00 00 00 00 00 00 00 00 00 00 00 00  ...D............
	0x0020 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 type."".NonEmptyInterface+0
	rel 8+8 t=1 type."".T+0
	rel 24+8 t=-32767 "".(*T).M1+0
	rel 32+8 t=-32767 "".(*T).M2+0
go.itab.*os.File,io.Writer SRODATA dupok size=32
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 44 b5 f3 33 00 00 00 00 00 00 00 00 00 00 00 00  D..3............
	rel 0+8 t=1 type.io.Writer+0
	rel 8+8 t=1 type.*os.File+0
	rel 24+8 t=-32767 os.(*File).Write+0
type..importpath.fmt. SRODATA dupok size=5
	0x0000 00 03 66 6d 74                                   ..fmt
gclocals·33cdeccccebe80329f1fdbee7f5874cb SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
"".T.M1.arginfo1 SRODATA static dupok size=11
	0x0000 fe 00 08 fe 08 08 10 08 fd fd ff                 ...........
"".T.M2.arginfo1 SRODATA static dupok size=11
	0x0000 fe 00 08 fe 08 08 10 08 fd fd ff                 ...........
gclocals·7d2d5fca80364273fb07d5820a76fef4 SRODATA dupok size=8
	0x0000 03 00 00 00 00 00 00 00                          ........
gclocals·5281ddd282d2ee46e64efba251339d2f SRODATA dupok size=14
	0x0000 03 00 00 00 0b 00 00 00 00 00 0a 00 05 00        ..............
"".main.stkobj SRODATA static size=80
	0x0000 03 00 00 00 00 00 00 00 c8 ff ff ff 10 00 00 00  ................
	0x0010 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0020 d8 ff ff ff 10 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 e8 ff ff ff 18 00 00 00  ................
	0x0040 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 24+8 t=1 runtime.gcbits.02+0
	rel 48+8 t=1 runtime.gcbits.02+0
	rel 72+8 t=1 runtime.gcbits.02+0
gclocals·09cf9819fc716118c209c2d2155a3632 SRODATA dupok size=10
	0x0000 02 00 00 00 02 00 00 00 02 00                    ..........
gclocals·69c1753bd5f81501d95132d08af04464 SRODATA dupok size=8
	0x0000 02 00 00 00 00 00 00 00                          ........
"".NonEmptyInterface.M1.arginfo1 SRODATA static dupok size=7
	0x0000 fe 00 08 08 08 fd ff                             .......
"".NonEmptyInterface.M2.arginfo1 SRODATA static dupok size=7
	0x0000 fe 00 08 08 08 fd ff                             .......
gclocals·dc9b0298814590ca3ffc3a889546fc8b SRODATA dupok size=10
	0x0000 02 00 00 00 02 00 00 00 03 00                    ..........
type..eq."".T.arginfo1 SRODATA static dupok size=5
	0x0000 00 08 08 08 ff                                   .....
gclocals·1a65e721a2ccc325b382662e7ffee780 SRODATA dupok size=10
	0x0000 02 00 00 00 01 00 00 00 01 00                    ..........
"".(*T).M1.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
"".(*T).M2.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
os.(*File).close.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
