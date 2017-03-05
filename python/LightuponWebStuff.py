import requests
import json

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
