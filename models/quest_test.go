package models

import (
  "testing"
  "github.com/davecgh/go-spew/spew"
  "fmt"
  "gopkg.in/yaml.v2"
  )

type Quest struct {
  Description string
  Posts []struct {
    Latitude string
    Longitude string
    Caption string
  }
}

type T struct {
        A string
        B struct {
                RenamedC int   `yaml:"c"`
                D        []int `yaml:",flow"`
        }
}

func TestParseQuestFromYaml(t *testing.T) {
  questYaml := `
description : "This is the dam tour."
posts:
  - latitude : 71.2345
    longitude : -42.734578
    caption : "Take all your dam pictures here."
  - latitude : 71.2534
    longitude : -42.75453
    caption : "This is the second dam post."
`
//   questYaml := `
// description : "Bar."
// posts:
// `


  

  var quest Quest
  err := yaml.Unmarshal([]byte(questYaml), &quest); if err != nil {
    fmt.Println("ERROR: Fuuuuuck that yaml.", err)
  }
  

  spew.Dump(quest)


}

func TestOtherThing (test *testing.T) {
  var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

  t := T{}
    
  err := yaml.Unmarshal([]byte(data), &t)
  if err != nil {
          fmt.Printf("error: %v", err)
  }
  fmt.Printf("--- t:\n%v\n\n", t)
}
