from bs4 import BeautifulSoup
import random
import requests
import json
import time
import LightuponWebStuff

NUM_SEARCH_PAGES = 10
# cityName = 'saltlake'
cityName = 'boston'

def curlAndSaveAtlasObscuraSearchPages():
	# boston
	searchLat = 42.358625
	searchLon = -71.1020555
	# saltlake
	searchLat = 40.772362
	searchLon = -111.889667
	for pageNum in range(1, NUM_SEARCH_PAGES):
		print pageNum
		searchURL = 'http://www.atlasobscura.com/search?lat=' + str(searchLat) + '&lng=' + str(searchLon) + '&q=&formatted_address=&source=desktop&nearby=false&page=' + str(pageNum)
		fileName = 'html/' + cityName + '/atlasobscura_' + cityName + '_page' + str(pageNum) + '.html'
		LightuponWebStuff.writeWebPageToFile(fileName, searchURL)
		time.sleep(random.randint(20,100)) ## to obfuscate the scraping, sleep a random amount of time

def readAndScrapeAtlasObscuraSearchPages():
	ping = False
	for pageNum in range(1, NUM_SEARCH_PAGES):
		print pageNum
		fileName = 'html/' + cityName + '/atlasobscura_' + cityName + '_page' + str(pageNum) + '.html'
		soup = LightuponWebStuff.getSoupFromFile(fileName)
		scenes = soup.find(id="js-search-results-list").find_all("div", class_="padding-container")


		for sceneDiv in scenes:

			try:
				uri = sceneDiv.find_all("a", class_="content-card-place")[0].attrs["href"]
				url = 'http://www.atlasobscura.com' + uri
				fileName = 'html/' + cityName + '/' + uri[8:] + '.html'
				# if (fileName == 'html/' + cityName + '/polcari-s-coffee.html'):
				# 	ping = True
				print 'fileName: ' + fileName

				
				# ## if we're gonna ping em, let's sleep for a second
				if (ping == True):
					print 'pinging...'
					time.sleep(random.randint(10,20)) ## to obfuscate the scraping, sleep a random amount of time
					LightuponWebStuff.writeWebPageToFile(fileName, url) #if it's our first time through, we'll need to uncomment this

				# Get most of the info from the search result page
				infoDiv = sceneDiv.find_all("div", class_="col-md-7")[0]
				title = str(infoDiv.find_all("span", class_="title-underline")[0].getText().encode("utf8"))
				description = str(infoDiv.find_all("div", class_="content-card-subtitle")[0].getText())
				imageURL = 'http://' + sceneDiv.find_all("div", class_="col-md-5")[0].find_all("img")[0].attrs["data-src"][2:]

				# Get the rest from the scene detail page
				sceneSoup = LightuponWebStuff.getSoupFromFile(fileName)
				latLonDivArray = sceneSoup.find("div", {"id": "lat-lng-element"}).find_all("div")
				latitude = float(latLonDivArray[0].getText().split('\n')[1] + '00')
				longitude = float(latLonDivArray[1].getText().split('\n')[1] + '00')

				# if (ping == True):
				LightuponWebStuff.postTrip(title, description, imageURL, latitude, longitude)
			# except (UnicodeEncodeError, IOError):
			except (UnicodeEncodeError):
				print 'failed to post trip: ' + title

def postAllThatWizardStuff():
	LightuponWebStuff.throwSomeDs(1, 'inventory', '(placeholder caption)', 'https://s-media-cache-ak0.pinimg.com/564x/ff/ac/b4/ffacb467ca843fa827896acf617f8b85.jpg')
	LightuponWebStuff.throwSomeDs(3, 'wizard', 'a wizard approaches...', 'https://s-media-cache-ak0.pinimg.com/originals/97/ce/ec/97ceec52d1c5b6aced484670a96e3e03.jpg')
	LightuponWebStuff.throwSomeDs(3, 'troll', 'find the troll...', 'https://s-media-cache-ak0.pinimg.com/originals/4f/b1/4f/4fb14fd09de912f4f2ac43b8aab93a3c.jpg')
	LightuponWebStuff.throwSomeDs(3, 'forest', 'the forest beckons...', 'https://s-media-cache-ak0.pinimg.com/originals/96/7c/c7/967cc7eee3c0c74ddd1a90edd39e8b7c.jpg')
	LightuponWebStuff.throwSomeDs(3, 'sword', 'you can do stuff with this', 'http://25.media.tumblr.com/tumblr_m8ud5a3Hoj1rrjmgoo1_1280.jpg')
	LightuponWebStuff.throwSomeDs(10, 'coin', 'one coin of the eastern realm', 'https://media1.britannica.com/eb-media/98/42598-004-02937021.jpg')
	LightuponWebStuff.throwSomeDs(10, 'ancient coin', 'one coin of the ancient realm', 'http://www.allaboutdrawings.com/image-files/japanese-sketch.jpg')

def postAllThatWizardStuff2():
	LightuponWebStuff.throwSomeDs(1, 'indian burial ground', 'this place is cursed as fuck','https://s-media-cache-ak0.pinimg.com/originals/d6/f5/12/d6f512ab0113931a23337c90979a861d.jpg')
	LightuponWebStuff.throwSomeDs(1, 'troll', 'find the troll', 'https://s-media-cache-ak0.pinimg.com/originals/4f/b1/4f/4fb14fd09de912f4f2ac43b8aab93a3c.jpg')
	LightuponWebStuff.throwSomeDs(1, 'wizard', 'a wizard approaches', 'https://s-media-cache-ak0.pinimg.com/originals/97/ce/ec/97ceec52d1c5b6aced484670a96e3e03.jpg')
	LightuponWebStuff.throwSomeDs(1, 'troll', 'find the troll', 'https://s-media-cache-ak0.pinimg.com/originals/4f/b1/4f/4fb14fd09de912f4f2ac43b8aab93a3c.jpg')
	LightuponWebStuff.throwSomeDs(1, 'wizard', 'a wizard approaches', 'https://s-media-cache-ak0.pinimg.com/originals/97/ce/ec/97ceec52d1c5b6aced484670a96e3e03.jpg')
	LightuponWebStuff.throwSomeDs(1, 'forest', 'this place is haunted as fuck', 'https://s-media-cache-ak0.pinimg.com/originals/96/7c/c7/967cc7eee3c0c74ddd1a90edd39e8b7c.jpg')
	LightuponWebStuff.throwSomeDs(1, 'indian burial ground', 'this place is cursed as fuck','https://s-media-cache-ak0.pinimg.com/originals/d6/f5/12/d6f512ab0113931a23337c90979a861d.jpg')
	LightuponWebStuff.throwSomeDs(1, 'wizard', 'a wizard approaches', 'https://s-media-cache-ak0.pinimg.com/originals/97/ce/ec/97ceec52d1c5b6aced484670a96e3e03.jpg')
	LightuponWebStuff.throwSomeDs(1, 'forest', 'this place is haunted as fuck', 'https://s-media-cache-ak0.pinimg.com/originals/96/7c/c7/967cc7eee3c0c74ddd1a90edd39e8b7c.jpg')
	LightuponWebStuff.throwSomeDs(1, 'indian burial ground', 'this place is cursed as fuck','https://s-media-cache-ak0.pinimg.com/originals/d6/f5/12/d6f512ab0113931a23337c90979a861d.jpg')
	LightuponWebStuff.throwSomeDs(1, 'troll', 'find the troll', 'https://s-media-cache-ak0.pinimg.com/originals/4f/b1/4f/4fb14fd09de912f4f2ac43b8aab93a3c.jpg')
	LightuponWebStuff.throwSomeDs(1, 'indian burial ground', 'this place is cursed as fuck','https://s-media-cache-ak0.pinimg.com/originals/d6/f5/12/d6f512ab0113931a23337c90979a861d.jpg')
	LightuponWebStuff.throwSomeDs(1, 'forest', 'this place is haunted as fuck', 'https://s-media-cache-ak0.pinimg.com/originals/96/7c/c7/967cc7eee3c0c74ddd1a90edd39e8b7c.jpg')




## BE FUCKING CAREFUL WITH THIS! It reaches out to atlas obscura and we don't want to make ourselves known
# curlAndSaveAtlasObscuraSearchPages()


# here's our main stuff here
readAndScrapeAtlasObscuraSearchPages()
# postAllThatWizardStuff2()
# LightuponWebStuff.postTrip('foster street', 'http://www.shorpy.com/node/3940', 'http://www.shorpy.com/files/images/boston.preview.jpg', 42.367533, -71.054213)
# LightuponWebStuff.postTrip('morrissey', 'http://www.vanyaland.com/2015/09/09/bostons-delta-bravo-urban-exploration-team-pinpoints-exact-location-of-80s-morrissey-photo/', 'http://www.vanyaland.com/wp-content/uploads/2015/05/MozBoston-1100x554.jpg', 42.297382, -71.048872)
# LightuponWebStuff.postTrip('coolidge corner', 'an old ass place', 'http://brooklinehistoricalsociety.org/archives/images/gs_21.jpg', 42.342011, -71.121342)
# LightuponWebStuff.postTrip('boston common', 'an old ass place', 'http://mfas3.s3.amazonaws.com/objects/SC157126.jpg', 42.355470, -71.066473)
# LightuponWebStuff.postTrip('boston massacre', '...known as the Incident on King Street by the British', 'http://cdn.history.com/sites/2/2015/04/hith-boston-massacre-152189046.jpg', 42.358815, -71.056629)
# LightuponWebStuff.postTrip('coolidge corner', 'this place is old as shit', 'http://brooklinehistoricalsociety.org/archives/images/coolidge_corner_1906.jpg', 42.342160, -71.121125)
# LightuponWebStuff.throwSomeDs(3, 'indian burial ground', 'this place is cursed as fuck','https://s-media-cache-ak0.pinimg.com/originals/d6/f5/12/d6f512ab0113931a23337c90979a861d.jpg')

