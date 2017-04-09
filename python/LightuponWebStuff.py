import requests
import json
import random
from bs4 import BeautifulSoup
import urllib

def getHttpStuff(uri, isDev):
	if (isDev):
		url = 'http://localhost:5000/lightupon/' + uri
		headers = {'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTc5NzEwNTUsImZhY2Vib29rSWQiOiIxMTQ1MjU2NDgyMTU0MDU1In0.ji2VyJmDuxiBnBYd19gGqvb7GzAGoVBf0lngZ3UzceA'}
	else:
		url = 'http://45.55.160.25/lightupon/' + uri
		headers = {'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI1MjEzMjYsImZhY2Vib29rSWQiOiIxMTQ1MjU2NDgyMTU0MDU1In0.hORIgWMv_fhjGbCQJK9hKCB63mk2gxlGrDClpyUgEfg'}
	return {'url' : url, 'headers' : headers}

def postTrip(title, description, imageURL, latitude, longitude):
	httpStuff = getHttpStuff('trips/generate', True)
	url = httpStuff['url']
	headers = httpStuff['headers']

	cards = [{'Caption': description, 'ImageURL' : imageURL}]
	payload = {'Latitude': latitude, 'Longitude' : longitude, 'Name': title, 'Route': 'aksdjhfkasjhdf', 'Street': '12', 'BackgroundURL':imageURL, 'Cards':cards}
	r = requests.post(url, data=json.dumps(payload), headers=headers)
	print 'Posting Trip: ', title, '    response: ', r

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

def getSoupFromFile(fileName):
	with open(fileName, 'r') as myfile:
	    data=myfile.read()
	return BeautifulSoup(data, "lxml")

def writeWebPageToFile(fileName, url):
	r = urllib.urlopen(url).read()
	text_file = open(fileName, "w")
	text_file.write(r)
	text_file.close()