package models

import (
  "testing"
  "github.com/davecgh/go-spew/spew"
  // "fmt"
  "gopkg.in/yaml.v2"
  )

type Quest struct {
  Foo string
}

func TestParseQuestFromYaml(t *testing.T) {
  questYaml := `
Quest:
  Description : "This is the dam tour."
  Posts:
    - Latitude : 71.2345
      Longitude : -42.734578
      Caption : "Take all your dam pictures here."
    - Latitude : 71.2534
      Longitude : -42.75453
      Caption : "This is the second dam post."
`
  var quest Quest
  err = yaml.Unmarshal(questYaml, &quest)

  spew.Dump(quest)

}
