import datetime


class EventDetailResponseModel:

    def __init__(self, mongoObj):
        self.name = mongoObj['name']
        self.timeStamp = mongoObj['timeStamp']
        self.daysLeft = calculateRemainingDays(mongoObj['timeStamp'])
        self.date = "11.10.2021"
        self.dateInTurkish = "11 Nisan 2021"
        self.type = mongoObj['type']
        self.attandees = mongoObj['attandees']
        self.location = mongoObj['location']
        self.distance = 10
        self.restrictions = mongoObj['restrictions']
        self.imageUrls = mongoObj['imageUrls']
        self.themes = mongoObj['themes']
        self.notes = mongoObj['notes']
        self.todayViewCount = mongoObj['todayViewCount']
        self.address = mongoObj['address']
        self.city = mongoObj['city']
        self.county = mongoObj['county']
        self.gender = mongoObj['gender']
        self.likePercentage = int(mongoObj['likeTotal'] / mongoObj['totalNumberOfRatings'])
        self.generalRating = mongoObj['generalRatingTotal'] / mongoObj['totalNumberOfRatings']
        self.locationRating = mongoObj['locationRatingTotal'] / mongoObj['totalNumberOfRatings']
        self.eventRating = mongoObj['eventRatingTotal'] / mongoObj['totalNumberOfRatings']
        self.comments = mongoObj['comments']
        self.favourited = mongoObj['favourited']


class EventBasicResponseModel:

    def __init__(self, mongoObj):
        self.id = mongoObj['_id']
        self.name = mongoObj['name']
        self.timeStamp = mongoObj['timeStamp']
        self.daysLeft = calculateRemainingDays(mongoObj['timeStamp'])
        self.dateInTurkish = "11 Nisan 2021"
        self.type = mongoObj['type']
        self.likePercentage = int(mongoObj['likeTotal'] / mongoObj['totalNumberOfRatings'])


class EventWithDistanceResponseModel:

    def __init__(self, mongoObj):
        self.id = mongoObj['_id']
        self.name = mongoObj['name']
        self.timeStamp = mongoObj['timeStamp']
        self.daysLeft = calculateRemainingDays(mongoObj['timeStamp'])
        self.dateInTurkish = "11 Nisan 2021"
        self.distance = 11  # mongoObj['distance']
        self.type = mongoObj['type']
        self.likePercentage = int(mongoObj['likeTotal'] / mongoObj['totalNumberOfRatings'])


def calculateRemainingDays(eventsTimeStamp):
    currentTimeStamp = datetime.datetime.now().timestamp()

    remainingDays = int((int(eventsTimeStamp) - currentTimeStamp) // 86400)

    return remainingDays
