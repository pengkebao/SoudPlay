package Sound

import (
	"C"
	"syscall"
	"unsafe"
)

const (
	SND_SYNC        = 0x0
	SND_ASYNC       = 0x1
	SND_NODEFAULT   = 0x2
	SND_MEMORY      = 0x4
	SND_LOOP        = 0x8
	SND_NOSTOP      = 0x10
	SND_NOWAIT      = 0x2000
	SND_ALIAS       = 0x10000
	SND_ALIAS_ID    = 0x110000
	SND_FILENAME    = 0x20000
	SND_RESOURCE    = 0x40004
	SND_PURGE       = 0x40
	SND_APPLICATION = 0x80
)

func Play(filePath string,sync bool) error {
	funInDllFile, err := syscall.LoadLibrary("Winmm.dll") // 调用的dll文件
	if err != nil {
		return err
	}
	defer syscall.FreeLibrary(funInDllFile)
	funName := "PlaySound"
	win32Fun, _ := syscall.GetProcAddress(syscall.Handle(funInDllFile), funName)
	Flags:=0
	if (sync){
		Flags = SND_FILENAME | SND_SYNC;
	}else{
		Flags = SND_FILENAME | SND_ASYNC;
	}
	file := C.CString(filePath) //转换成char*
	callFun:=uintptr(win32Fun) //方法名
	paraNum:=uintptr(3) //方法参数个数
	para1:=uintptr(unsafe.Pointer(file)) //方法参数1
	para2:=uintptr(0)//方法参数2
	para3:=uintptr(Flags)//方法参数3
	r, _, errno := syscall.Syscall(callFun,paraNum,para1,para2,para3)
	if r==0{
		return errno
	}
	return nil
}
