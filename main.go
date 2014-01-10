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


const Tmpl = `---
title: {{title}}
author: {{author}}
image: {{image}}
tags: 
{{#tags}}
 - {{.}} {{! This wrong seeming setup means there won't be a extra newline between the tag list and the date }}
{{/tags}}date: {{date}}
---

{{body}}`


func input(field string, scanner *bufio.Scanner) string {
  fmt.Printf("%s:", field)
  scanner.Scan()
  return scanner.Text()
}

func collectData(scanner *bufio.Scanner) map[string]interface{} {
  metadata := make(map[string]interface{})
  metadata["title"] = input("Title", scanner)
  metadata["author"] = input("Author", scanner)
  metadata["image"] = input("Image path", scanner)
  metadata["tags"]= strings.Split(input("Tags", scanner),", ")
  return metadata
}

func getBody(scanner *bufio.Scanner) string {
  tmpFile := "/tmp/coconut-post"
  editor := os.Getenv("EDITOR")
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
  var afName string
  if len(os.Args) > 1 {
      afName = os.Args[1]
  }
  scanner := bufio.NewScanner(os.Stdin)
  article := collectData(scanner)
  article["body"] = getBody(scanner)
  if article["body"] == "" {
    fmt.Println("Aborting due to empty body")
    return
  }
  if afName == ""{
    afName = strings.Replace(strings.ToLower(article["title"].(string)), " ", "-", -1)
    if len(afName)> 20{
        afName=afName[:20]
    }
  }
  ctime := time.Now()
  article["date"] = ctime.Format(DateFormat)
  rendered := mustache.Render(Tmpl, article)
  aFile, err := os.Create(afName+".md")
  defer aFile.Close()
  if err != nil{
    fmt.Println("Error:", err)
    return
  }
  aFile.WriteString(rendered)
}
