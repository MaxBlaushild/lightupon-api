import requests
import json
import random

def postTrip(title, description, imageURL, latitude, longitude):
	# DEV - don't use my token tho ya punk bitch
	# url = 'http://localhost:5000/lightupon/trips/generate'
	# headers = {'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTc5NzEwNTUsImZhY2Vib29rSWQiOiIxMTQ1MjU2NDgyMTU0MDU1In0.ji2VyJmDuxiBnBYd19gGqvb7GzAGoVBf0lngZ3UzceA'}
	# PROD
	url = 'http://45.55.160.25/lightupon/trips/generate'
	headers = {'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI1MjEzMjYsImZhY2Vib29rSWQiOiIxMTQ1MjU2NDgyMTU0MDU1In0.hORIgWMv_fhjGbCQJK9hKCB63mk2gxlGrDClpyUgEfg'}
	cards = [{'Caption': description, 'ImageURL' : imageURL}]
	payload = {'Latitude': latitude, 'Longitude' : longitude, 'Name': title, 'Route': 'aksdjhfkasjhdf', 'Street': '12', 'BackgroundURL':imageURL, 'Cards':cards}
	r = requests.post(url, data=json.dumps(payload), headers=headers)

def throwSomeDs(numberOfDs, title, description, imageURL):
	for i in range(numberOfDs):
		postTripWithRandomLocation(title, description, imageURL)

def postTripWithRandomLocation(title, description, imageURL):
	latitude = randomLocation()[0]
	longitude = randomLocation()[1]
	postTrip(title, description, imageURL, latitude, longitude)

def randomLocation():
	# smaller area
	lat = random.uniform(42.332856, 42.342373)
	lon = random.uniform(-71.090436, -71.065889)
	# small area
	# lat = random.uniform(42.326854, 42.368061)
	# lon = random.uniform(-71.110049, -71.027308)
	# large area
	# lat = random.uniform(42.289632, 42.446392)
	# lon = random.uniform(-71.268320, -70.878992)
	return [lat, lon]