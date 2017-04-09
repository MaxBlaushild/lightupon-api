import LightuponWebStuff
import requests
import json


print 2+2

def getUserLocations():
	return {'1':54, '2':76}


# pull locations out of the DB 
locations = getUserLocations()
print locations

# run stats

# serialize result
payload = '[4,1,6]'
payload = '{"foo":"bar"}'

# post to endpoint
httpStuff = LightuponWebStuff.getHttpStuff('engagementScores', True)
url = httpStuff['url']
headers = httpStuff['headers']

print url
print headers

requests.post(url, data=json.dumps(payload), headers=headers)

