package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "io/ioutil"
    "crypto/sha256"
    "flag"
    "encoding/hex"
    "hash"
)

const BlockSize = 64

type item struct {
    path string
    parent *item
    children []*item
    hash  []byte
}

func print(item *item, offset int){
    fmt.Println(strings.Repeat("   ", offset), "|")
    fmt.Println(strings.Repeat("   ", offset), strings.Repeat("---", offset),
        filepath.Base(item.path), hex.EncodeToString(item.hash))
    for _, child := range item.children{
        print(child, offset+1)
    }
}

func getHash(item *item){
    hasher := sha256.New()
    s, err := ioutil.ReadFile(item.path)
    if err != nil {
        return
    }
    hasher.Write(s)
    item.hash = hasher.Sum(nil)
}

func folderHash(item *item, hasher hash.Hash){
    for _, child := range item.children {
        hasher.Write(child.hash)
        folderHash(child, hasher)
    }
}

func main() {
    dir := flag.String("dir", "testdir", "path to folder")
    count := flag.Int("count", 4, "run program in concurent mode with n goroutines")
    flag.Parse()
    root, err := filepath.Abs(*dir)
    if err != nil {
        log.Fatal(err)
    }
    folders := make(map[string]*item)
    saveItem := func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        folders[path] = &item{
            path: path,
            children: make([]*item, 0),
        }
        return nil
    }
    if err = filepath.Walk(root, saveItem); err != nil {
        log.Fatal(err)
    }
    var rootItem *item
    for path, item := range folders {
        parentItem, exists := folders[filepath.Dir(path)]
        if exists {
            item.parent = parentItem
            parentItem.children = append(parentItem.children, item)
        } else {
            rootItem = item
        }
    }
    ch := make (chan *item)
    var wg sync.WaitGroup
    wg.Add(*count)
    for i:=0; i<(*count); i++ {
        go func() {
            defer wg.Done()
            for {
                if item, more := <- ch; more {
                    getHash(item)
                } else {
                    return
                }
            }
        }()
    }
    for _, item := range folders {
        ch <- item
    }
    for _, item := range folders {
        if len(item.children) > 0 {
            hasher := sha256.New()
            folderHash(item, hasher)
            item.hash = hasher.Sum(nil)
        }
    }
    close(ch)
    wg.Wait()
    print(rootItem, 0)
}