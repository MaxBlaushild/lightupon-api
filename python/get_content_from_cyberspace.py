# import urllib2
# from xml.etree import ElementTree as etree
import lightupon_data_access_functions

rss_urls = [
	'http://rss.nytimes.com/services/xml/rss/nyt/Technology.xml',
	'http://motherboard.vice.com/rss?trk_source=motherboard',
	'http://www.reddit.com/r/videos/top/.rss'
]

for url in rss_urls:
	print url
	lightupon_data_access_functions.get_stuff_from_rss_url(url)