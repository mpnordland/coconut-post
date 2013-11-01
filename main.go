package main

import (
  "bufio"
  "os"
  "io/ioutil"
  "os/exec"
  "fmt"
  "time"
  "strings"
  "github.com/hoisie/mustache"
)

const DateFormat = "Jan _2 2006 15:04"


func input(field string, scanner *bufio.Scanner) string {
  fmt.Printf("%s:", field)
  scanner.Scan()
  return scanner.Text()
}

func collectData(scanner *bufio.Scanner) map[string]interface{} {
  metadata := make(map[string]interface{})
  metadata["title"] = input("Title", scanner)
  metadata["author"] = input("Author", scanner)
  metadata["tags"]= strings.Split(input("Tags", scanner),", ")
  return metadata
}

func getBody(scanner *bufio.Scanner) string {
  tmpFile := "/tmp/coconut-post"
  editor := "/usr/bin/vim" //strings.Split(os.Getenv("EDITOR"), " ")[0]
  if editor == "" {
    fmt.Println("What editor would you like to write your post in?")
    scanner.Scan()
    editor = scanner.Text()
  }
  cmd := exec.Command(editor, tmpFile)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  err := cmd.Run()
  if err == nil {
    cont, e := ioutil.ReadFile(tmpFile)
    defer os.Remove(tmpFile)
    if e == nil || len(cont) != 0 {
      return string(cont)
    }
  }
  return ""
}

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  article := collectData(scanner)
  article["body"] = getBody(scanner)
  if article["body"] == "" {
    fmt.Println("Aborting due to empty body")
    return
  }
  afName := strings.ToLower(article["title"].(string))
  if len(afName)> 20{
    afName=afName[:20]
  }
  ctime := time.Now()
  article["date"] = ctime.Format(DateFormat)
  rendered := mustache.RenderFile("post.tmpl", article)
  aFile, err := os.Create(afName+".md")
  defer aFile.Close()
  if err != nil{
    fmt.Println("Error:", err)
    return
  }
  aFile.WriteString(rendered)
}