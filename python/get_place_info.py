


from googleplaces import GooglePlaces, types, lang
import json
import urllib2


# from sklearn import feature_extraction as homie
from sklearn.feature_extraction.text import DictVectorizer




def get_details_for_place_id(place_id):
    print('place_id')
    print(place_id)
    url = 'https://maps.googleapis.com/maps/api/place/details/json?placeid='+place_id+'&key=AIzaSyAhsrriQNAkhtkCdJDN2nC1TOTjflbnYYg'
    data = json.load(urllib2.urlopen(url))
    print data['result']

# get_details_for_place_id('ChIJN1t_tDeuEmsRUsoyG83frY4')

# documents = [{'sup bro','hey whats on that flipside jack'}]
measurements = [
     {'city': 'Dubai', 'temperature': 33.},
     {'city': 'London', 'temperature': 12.},
     {'city': 'San Fransisco', 'temperature': 18.},
 ]
vec = DictVectorizer()

print vec.fit_transform(measurements).toarray()



# # documents = [open(f) for f in text_files]

# tfidf = TfidfVectorizer().fit_transform(documents)
# # no need to normalize, since Vectorizer will return normalized tf-idf
# pairwise_similarity = tfidf * tfidf.T
# print 'pairwise_similarity'
# print pairwise_similarity






print('_______________________________________________________________')
exit()

YOUR_API_KEY = 'AIzaSyAhsrriQNAkhtkCdJDN2nC1TOTjflbnYYg'

google_places = GooglePlaces(YOUR_API_KEY)



# You may prefer to use the text_search API, instead.
query_result = google_places.nearby_search(
        location='London, England', keyword='Fish and Chips',
        radius=200, types=[types.TYPE_FOOD])

# query_result

# print getattr(query_result)

# print query_result.__dict__

if query_result.has_attributions:
    print query_result.html_attributions


for place in query_result.places:
    # Returned places from a query are place summaries.
    print "--------------------------------------------"
    print place.id
    print place.name
    print place.geo_location
    # print place.place_id

    # # The following method has to make a further API call.
    place.get_details()
    # # Referencing any of the attributes below, prior to making a call to
    # # get_details() will raise a googleplaces.GooglePlacesAttributeError.
    # print place.details # A dict matching the JSON response from Google.
    print place.local_phone_number
    # print place.international_phone_number
    # print place.website
    print place.url

    print place.rating

    # for details in place.details:


    # Getting place photos

    # for photo in place.photos:
    #     # 'maxheight' or 'maxwidth' is required
    #     photo.get(maxheight=500, maxwidth=500)
    #     # MIME-type, e.g. 'image/jpeg'
    #     photo.mimetype
    #     # Image URL
    #     photo.url
    #     # Original filename (optional)
    #     photo.filename
    #     # Raw image data
    #     photo.data


