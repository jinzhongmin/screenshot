# screenshot 截图小工具
用go-gtk实现的截图小工具，在msys2下做的
## 安装
### msys2
1.安装 go

2安装go-gtk 和 screenshot

```  bash
pacman -S mingw-w64-x86_64-gtk2
go get github.com/mattn/go-gtk/gtk
go get github.com/vova616/screenshot
```

3.编译
``` bash
git clone https://github.com/jinzhongmin/screenshot
cd screenshot
go build
```

