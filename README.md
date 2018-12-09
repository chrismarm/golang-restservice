# Golang RESTful API

As developed in https://github.com/chrismarm/python-restservice with Python, in this project we show some features of Golang language. 

* The API offers methods to list, get, add, update and delete dive locations, that only have a name, depth, latitude and longitude. An UUID is automatically generated for new locations. 'Name' attribute is used to index locations, so they can be deleted and retrieved by this value, that is, methods guarantee that there will not be 2 locations with the same name.
* A JSON file is used as database to show file operations and JSON management. An array and a map are used as a cache for service methods and quicker retrieval. These collections must be correctly updated in some service methods.
* [Google Maps Geocoding API](https://developers.google.com/maps/documentation/geocoding/start) is accessed to get some extra information about the location defined by latitude/longitude. A key given by Google is needed to call this API, that is passed as the first command line argument. The information retrieved for a latitude/longitude pair are locality, country, a friendly description of the location and the [Plus Code](https://plus.codes/). This info is added and stored for every new inserted dive location.
* As the JSON response from Google Maps API is quite large, a utility called [GJSON](https://github.com/tidwall/gjson) has been used in order to extract the desired fields easily by using dot notation and iteration methods.

Future improvements:
* Unit tests
* Geocoding allows a wide variety of improvements, simply extracting information from latitude/longitude