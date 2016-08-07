package models

import(
      "log"
      "time"
      "github.com/jinzhu/gorm"
      "github.com/aws/aws-sdk-go/aws"
      "github.com/aws/aws-sdk-go/aws/session"
      "github.com/aws/aws-sdk-go/service/s3"
      )

type Scene struct {
  gorm.Model
  Name string
  Latitude float64
  Longitude float64
  TripID uint
  BackgroundUrl string
  SceneOrder uint
  Featured bool
  Cards []Card
  SoundKey string
  SoundResource string
}

func ShiftScenesUp(sceneOrder int, tripID int) bool {
  scene := Scene{}
  DB.Where("trip_id = $1 AND scene_order = $2", tripID, sceneOrder).First(&scene)
  if scene.ID == 0 {
    return true
  } else {
    ShiftScenesUp(sceneOrder + 1, 1)
    DB.Model(&scene).Update("scene_order", sceneOrder + 1)
    return true
  }
}

func ShiftScenesDown(sceneOrder int, tripID int) bool {
  scene := Scene{}
  DB.Where("trip_id = $1 AND scene_order = $2", tripID, sceneOrder + 1).First(&scene)
  if scene.ID == 0 {
    return true
  } else {
    ShiftScenesDown(sceneOrder + 1, 1)
    DB.Model(&scene).Update("scene_order", sceneOrder)
    return true
  }
}

func (s *Scene) PopulateSound() {
  svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
  req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
    Bucket: aws.String("lightupon"),
    Key:    aws.String(s.SoundKey),
  })

  urlStr, err := req.Presign(15 * time.Minute)

  if err != nil {
      log.Println("Failed to sign request", err)
  }

  s.SoundResource = urlStr
}