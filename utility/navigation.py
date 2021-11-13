import logging
import reverse_geocoder as rg
from geopy import distance


def getCity(x, y):
    coordinates = (x, y)

    logging.debug("Coordinates of the client, x: " + str(x) + " y: " + str(y))
    results = rg.search(coordinates)  # default mode = 2

    city = (list(results[0].items())[3][1])  # dont
    logging.debug(city)

    return city


def calculateDistance(x1, y1, x2, y2):
    distance_between = distance.distance((x1, y1), (x2, y2)).km
    return distance_between
