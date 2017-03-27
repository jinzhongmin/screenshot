# screenshot
screenshot tool with go-gtk
## How to build
### msys2
1.installed go

2.install go-gtk and screenshot

```  bash
pacman -S mingw-w64-x86_64-gtk2
go get github.com/mattn/go-gtk/gtk
go get github.com/vova616/screenshot
```

3.build it

``` bash
git clone github.com/jinzhongmin/screenshot
cd screenshot
go build
```

