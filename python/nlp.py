import nltk
import string
import os

import json
import urllib2

from sklearn.feature_extraction.text import TfidfVectorizer
from nltk.stem.porter import PorterStemmer

from sklearn.feature_extraction import DictVectorizer






path = '/opt/datacourse/data/parts'
token_dict = {}
stemmer = PorterStemmer()



########################################################################################################
import unicodedata
import sys

tbl = dict.fromkeys(i for i in xrange(sys.maxunicode)
                      if unicodedata.category(unichr(i)).startswith('P'))
def remove_punctuation(text):
    return text.translate(tbl)
########################################################################################################



def stem_tokens(tokens, stemmer):
    stemmed = []
    for item in tokens:
        stemmed.append(stemmer.stem(item))
    return stemmed

def tokenize(text):
    tokens = nltk.word_tokenize(text)
    stems = stem_tokens(tokens, stemmer)
    return stems


YOUR_API_KEY = 'AIzaSyAhsrriQNAkhtkCdJDN2nC1TOTjflbnYYg'

def get_reviews_for_place_id(place_id):
    print('place_id')
    print(place_id)
    url = 'https://maps.googleapis.com/maps/api/place/details/json?placeid='+place_id+'&key=AIzaSyAhsrriQNAkhtkCdJDN2nC1TOTjflbnYYg'
    data = json.load(urllib2.urlopen(url))
    # print data['result']['reviews']
    reviews = data['result']['reviews']
    token_dict = {}
    for review in reviews:
        text = review['text']
        lowers = text.lower()
        no_punctuation = remove_punctuation(lowers)
        token_dict[counter] = no_punctuation
    return token_dict
	



def search_for_thing_at_location(thing, latitude, longitude):
    url = 'https://maps.googleapis.com/maps/api/place/nearbysearch/json?keyword='+thing+'&key=AIzaSyAhsrriQNAkhtkCdJDN2nC1TOTjflbnYYg'
    print url
    # data = json.load(urllib2.urlopen(url))
    # print data['result']['reviews']
    # reviews = data['result']['reviews']
    # token_dict = {}
    # for review in reviews:
    #     text = review['text']
    #     lowers = text.lower()
    #     no_punctuation = remove_punctuation(lowers)
    #     token_dict[counter] = no_punctuation
    # return token_dict
    return 78;

print search_for_thing_at_location('historical', 42.347377, -71.119276)


# reviews = get_reviews_for_place_id('ChIJN1t_tDeuEmsRUsoyG83frY4')
# print(reviews)
# counter = 0
# for review in reviews:
# 	counter += 1
# 	text = review['text']
# 	lowers = text.lower()
# 	# print type(lowers)
# 	no_punctuation = remove_punctuation(lowers)
# 	token_dict[counter] = no_punctuation





# greetings = {}
# greetings[0] = 'cake cake cake cake dogs'
# greetings[1] = 'cake cake cake cake cats'

#this can take some time
# tfidf = TfidfVectorizer(tokenizer=tokenize, stop_words='english')
# tfs = tfidf.fit_transform(token_dict.values())
# tfs = tfidf.fit_transform(greetings.values())
# print tfs.toarray()

# print tfs.get_feature_names()
# feature_names = tfidf.get_feature_names()
# print feature_names


# for i in range(0,len(feature_names)):
# 	print feature_names[i]
# 	print key, value




# vec = DictVectorizer()
# print vec.fit_transform(token_dict.values()).toarray()






