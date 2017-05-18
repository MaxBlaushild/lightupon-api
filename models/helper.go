package models

import("math")

func CalculateDistance(location1 UserLocation, location2 UserLocation) (distance float64){
  R := 6371000.0
  phi1 := toRadians(location1.Latitude)
  phi2 := toRadians(location2.Latitude)
  deltaLambda := toRadians(location2.Longitude - location1.Longitude)
  distance = math.Acos(math.Sin(phi1)*math.Sin(phi2) + math.Cos(phi1)*math.Cos(phi2) * math.Cos(deltaLambda) ) * R
  return
}

func toRadians(angleInDegrees float64) (angleInRadians float64) {
  angleInRadians = angleInDegrees*(3.145/180.001)
  return
}

// TODO: this function needs to be destroyed, along with refactoring to not have Location and UserLocation as separate things
func CalculateLocationDistance(location1 Location, location2 Location) (distance float64){
  userLocation1 := UserLocation{Latitude : location1.Latitude, Longitude : location1.Longitude}
  userLocation2 := UserLocation{Latitude : location2.Latitude, Longitude : location2.Longitude}
  distance = CalculateDistance(userLocation1, userLocation2)
  return
}