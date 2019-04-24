package main;
//  递归监控会导致无法删除和修改父目录
// 一般不需要递归监控
import (
    "github.com/fsnotify/fsnotify"
    "fmt"
    "path/filepath"
    "os"
)
 
type Watch struct {
    watch *fsnotify.Watcher;
}
var WatchPaths []string
//监控目录
func (w *Watch)DirectoryTraversal(dir string){
	for i:=0;i<len(WatchPaths);i++{
		tmpp:=WatchPaths[i]
		err := w.watch.Remove(tmpp)
		if err != nil {
			fmt.Println("监控删除 err : ", err)
		}else{
			fmt.Println("监控删除 : ", tmpp)
		}
	}
	WatchPaths=WatchPaths[0:0]
	//通过Walk来遍历目录下的所有子目录
    filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        //这里判断是否为目录，只需监控目录即可
        //目录下的文件也在监控范围内，不需要我们一个一个加
        if info.IsDir() {
            path, err := filepath.Abs(path);
            if err != nil {
                return err
			}
			fmt.Println("add path : ", path)
			WatchPaths=append(WatchPaths,path)
        }
        return nil;
	});
	for i:=0;i<len(WatchPaths);i++{
		tmpp:=WatchPaths[i]
		err := w.watch.Add(tmpp)
		if err != nil {
			return 
		}
		fmt.Println("监控 : ", tmpp)
	}
}
func (w *Watch) watchDir(dir string) {
	w.DirectoryTraversal(dir)
    go func() {
        for {
            select {
            case ev := <-w.watch.Events:
                {
                    if ev.Op&fsnotify.Create == fsnotify.Create {
                        fmt.Println("创建文件 : ", ev.Name);
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						w.DirectoryTraversal(dir)
                        // fi, err := os.Stat(ev.Name);
                        // if err == nil && fi.IsDir() {
						// 	w.DirectoryTraversal(dir)
                        // }
                    }
                    if ev.Op&fsnotify.Write == fsnotify.Write {
                        fmt.Println("写入文件 : ", ev.Name);
                    }
                    if ev.Op&fsnotify.Remove == fsnotify.Remove {
                        fmt.Println("删除文件 : ", ev.Name);
                        //如果删除文件是目录，则移除监控
                        fi, err := os.Stat(ev.Name);
                        if err == nil && fi.IsDir() {
							w.DirectoryTraversal(dir)
                        }
                    }
                    if ev.Op&fsnotify.Rename == fsnotify.Rename {
                        fmt.Println("重命名文件 : ", ev.Name);
                    }
                    if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
                        fmt.Println("修改权限 : ", ev.Name);
                    }
                }
            case err := <-w.watch.Errors:
                {
                    fmt.Println("error : ", err);
                    return;
                }
            }
        }
    }();
}
 
func main() {
    watch, _ := fsnotify.NewWatcher()
    w := Watch{
        watch: watch,
    }
    w.watchDir(`G:\test`);
    select {};
}