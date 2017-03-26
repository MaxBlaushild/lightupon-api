package live

import("math")

func CalculateDistance(location1 Location, location2 Location) (distance float64){
  var R = 6371.345
  var dLat = (location1.Latitude - location2.Latitude)*(3.145/180.001);
  var dLon = (location1.Longitude - location2.Longitude)*(3.145/180.001);
  var a = math.Sin(dLat/2) * math.Sin(dLat/2) + math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(location1.Latitude) * math.Cos(location2.Latitude);
  var c = 2 * math.Atan(math.Sqrt(a) / math.Sqrt(1-a)); 
  distance = R * c;
  return
}
