package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "flag"
    "net/http"
    "io"
)

func panic(err error) {
    if err != nil {
        fmt.Println(err)
        panic(err)
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

func c_parse(url string, i int, c chan int64){
    size := get_file(url, i)
    c <- size
}

func parse(url string, i int){
    fmt.Println(get_file(url, i))
}

func main() {
    filePtr := flag.String("file_path", "test.txt", "path to file with url")
    concurentPtr := flag.Bool("concurent", false, "run program in concurent mode")
    flag.Parse()
    fmt.Println(*concurentPtr)
    file, err := os.Open(*filePtr)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    c := make(chan int64)
    defer close(c)
    i := 0
    for scanner.Scan() {
        if *concurentPtr {
            go c_parse(scanner.Text(), i, c)
        } else {
            parse(scanner.Text(), i)
        }
        i++
    }
    if *concurentPtr {
        for j:=0; j<i; j++ {
            fmt.Println(<-c)
        }
    }
}
