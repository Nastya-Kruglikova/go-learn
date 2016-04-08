package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "flag"
    "net/http"
    "io"
    "sync"
)

func panic(err error) {
    if err != nil {
        fmt.Println(err)
        log.Fatal(err)
    }
}

func get_file(url string, i int) int64{
    fmt.Println(url)
    resp, err := http.Get(url)
    defer resp.Body.Close()
    panic(err)
    file, err := os.Create(fmt.Sprintf("res%d.txt", i))
    panic(err)
    defer file.Close()
    fmt.Println(resp.Status)
    size, err := io.Copy(file, resp.Body)
    panic(err)
    return size
}

func main() {
    filePtr := flag.String("file_path", "test.txt", "path to file with url")
    concurentPtr := flag.Int("concurent", 4, "run program in concurent mode with n goroutines")
    var wg sync.WaitGroup
    flag.Parse()
    file, err := os.Open(*filePtr)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    wg.Add(*concurentPtr)
    for i:=0; i<(*concurentPtr); i++ {
        go func() {
            defer wg.Done()
            for scanner.Scan() {
                get_file(scanner.Text(), i)
            }
        }()
    }
    wg.Wait()
}
