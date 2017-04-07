import psycopg2
import os
import urllib2
from xml.etree import ElementTree as etree


## Connect to the database

# def initialize():
db_name = os.environ['DB_NAME']
db_username = os.environ['DB_USERNAME']
# conn = psycopg2.connect("dbname=" + db_name + " user=" + db_username)

def open_connection():
	return psycopg2.connect("dbname=" + db_name + " user=" + db_username)


def save_bookmark(title, url, pubDate):
	conn = open_connection(); cur = conn.cursor()
	cur.execute("INSERT INTO bookmarks (title, url, pub_date) VALUES (%s, %s, %s)", (title, url, pubDate))
	conn.commit(); cur.close(); conn.close()
	

# Really only gets one scene
def getOneScene():
	conn = open_connection()
	cur = conn.cursor()
	cur.execute("SELECT * FROM scenes;")
	return cur.fetchone()


def get_stuff_from_rss_url(rss_url_string):
	rss_file = urllib2.urlopen(rss_url_string)

	#convert to string:
	rss_data = rss_file.read()

	#close file because we dont need it anymore:
	rss_file.close()

	#entire feed
	reddit_root = etree.fromstring(rss_data)
	items = reddit_root.findall('channel/item')

	title = ''; url = ''

	for item in items:   
		for child in item:
			if (child.tag == 'title'):
				title = child.text
			elif (child.tag == 'link'):
				url = child.text
			elif (child.tag == 'pubDate'):
				pubDate = child.text
				
			
		if (title and url and pubDate):
			print '  ', title, url
			save_bookmark(title, url, pubDate)
		else:
			print 'ERROR ERROR ERROR ERROR ERROR ERROR ERROR ERROR ERROR ERROR ERROR ERROR ERROR ERROR'

def do_things():
	print 'stuff and things'
