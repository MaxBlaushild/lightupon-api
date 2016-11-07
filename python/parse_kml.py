# import untangle
# obj = untangle.parse('FreedomTrail.xml')
# print type(obj.root)





import xmltodict

with open('FreedomTrail.xml') as fd:
    doc = xmltodict.parse(fd.read())

kml = doc['kml']['Document']['Folder']['Placemark']

print len(kml)
for place in kml:
    print place.keys()
    

# print kml.keys()