package main

import (
    "github.com/jlaffaye/ftp"
    "github.com/akamensky/argparse"
    "io/ioutil"
    "fmt"
    "log"
    "os"
)


var default_server string = "vm:21"
var default_user string = "test"
var default_password string = "test"
var default_directory = "test"


func check(err error){
    if err != nil {
    panic(err)
    }
}


func put(remote_dirname *string, filename *string, c *ftp.ServerConn) {

    r, err := os.Open(*filename)
    check(err)

    if len(*remote_dirname) == 0{
        *remote_dirname = default_directory
    }
    c.ChangeDir(*remote_dirname)
    // 'dir/file' 553
    err = c.Stor(*filename, r)
    check(err)
}


func get(remote_dirname *string, filename *string, c *ftp.ServerConn) {

    remote_file := conf_remote_file(remote_dirname, filename)

    r, err := c.Retr(remote_file)
    check(err)
    buf, err := ioutil.ReadAll(r)
    check(err)
    ioutil.WriteFile(*filename, buf, 0644)
}


func list(remote_dirname *string, filename *string, c *ftp.ServerConn) {

    remote_file := conf_remote_file(remote_dirname, filename)
    entries, _ := c.List(remote_file)

    for _, entry := range entries {
        name := entry.Name
        fmt.Println(name)
    }
}


func conf_server_auth(s *string, u *string, p *string){
    if *s == ""{*s=default_server}
    if *u == ""{*u=default_user}
    if *p == ""{*p=default_password}
}

func conf_remote_file(remote_dirname *string, filename *string)(string) {
    // string "dir/file" representing abs file path on the server
    // only available for ftp {get/list} method but NOT for ftp {put} method
    if len(*remote_dirname) == 0{
        *remote_dirname = default_directory
    }
    remote_file := fmt.Sprintf("%s/%s", *remote_dirname, *filename)
    return remote_file
}


func login_ftp(s string, u string, p string)(*ftp.ServerConn) {
    c, err := ftp.Dial(s)   // c: client or connection
    if err != nil {
    log.Fatal(err)
    }

    if err := c.Login(u, p); err != nil {
    log.Fatal(err)
    }

    return c
}


func main(){

    parser := argparse.NewParser("FTP Client Powered by Yao", "WoW, this is amazing and there you go!")
    // FTP server info
    server := parser.String("s", "server", &argparse.Options{Help: "FTP server", Required: false })
    user := parser.String("u", "user", &argparse.Options{Help: "FTP server user", Required: false })
    password := parser.String("p", "password", &argparse.Options{Help: "FTP server user password", Required: false })
    // FTP method
    method := parser.String("m", "method", &argparse.Options{Required: false, Help: "FTP {method}",  Default: "list"})
    // FTP server side
    remote_dirname := parser.String("d", "dir", &argparse.Options{Help: "FTP server dir", Required: false })

    filename := parser.String("f", "filename", &argparse.Options{Help: "FTP server file inside dir", Required: false })

    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Print(parser.Usage(err))
    }else{

        conf_server_auth(server, user, password)

        c := login_ftp(*server, *user, *password)

        switch *method{
        case "p", "put":
            put(remote_dirname, filename, c)
        case "g", "get":
            get(remote_dirname, filename, c)
        default:
            list(remote_dirname, filename, c)
        }
    }
}

