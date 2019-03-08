import sys
from pykml import parser
from yaml import dump, Dumper

"""
This here script is for transforming kml files into quests. To use it, download a kml file, find yourself a nice image 
to use for the quest and grab the image url, decide which user id and quest id to use (for now, this quest must exist already), and run the following:

	python kml2quest.py {filename} {image_url} {quest_id} {user_id}

e.g.

	python kml2quest.py example_kml_files/chinatown_san_francisco.kml https://pixfeeds.com/images/usa-travel/1280-478192484-gateway-arch-on-grant-avenue.jpg 8 5

(You can find many cool kml's by searching "site:maps.google.com cool maps search terms" and downloading the kml).

"""

print '__________________ PROCESSING KML... __________________'

# ________________________ INPUTS ________________________

filename = sys.argv[1]
imageurl = sys.argv[2]
questid = int(sys.argv[3]) # TODO: nicely handle invalid input that fails to be parsed to an integer
userid = int(sys.argv[4])

file = open(filename, 'r')
root = parser.fromstring(file.read())
file.close()

# ________________________ MAIN ALGO ________________________

quest = {}
quest['description'] = str(root.Document.name)
quest['timetocomplete'] = 0
quest['userid'] = userid

posts = []

for placemark in root.Document.Folder.Placemark:

	post = {}
	post['caption'] = str(placemark['name']) ## can also use key "description"
	post['questid'] = questid
	post['imageurl'] = imageurl

	try:
		splert = str(placemark.Point.coordinates).replace('\n', '').replace(' ', '').split(',')
	except:
		print 'ERROR: Couldnt get the location for placemark. Lets just skip it.'
		continue

	# Longitude goes first because kml is backwards. Whyyyyyyy...
	post['longitude'] = float(splert[0])
	post['latitude'] = float(splert[1])
	
	posts.append(post)

quest['posts'] = posts

print '__________________ OK HERES THE QUEST YAML __________________\n'

print dump(quest, default_flow_style=False)
