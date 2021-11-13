from app import mongo_events


class EventModel():
    """
    Event Model Class
    """

    def __init__(self, mongoObj):
        self.name = mongoObj['name']
        self.timeStamp = mongoObj['timeStamp']
        self.type = mongoObj['type']
        self.attandees = mongoObj['attandees']
        self.location = mongoObj['location']
        self.organizer = mongoObj['organizer']
        self.coordinates = mongoObj['coordinates']
        self.restrictions = mongoObj['restrictions']
        self.imageUrls = mongoObj['imageUrls']
        self.themes = mongoObj['themes']
        self.notes = mongoObj['notes']
        self.today = mongoObj['today']
        self.todayViewCount = mongoObj['todayViewCount']
        self.address = mongoObj['address']
        self.city = mongoObj['city']
        self.county = mongoObj['county']
        self.gender = mongoObj['gender']
        self.totalNumberOfRatings = mongoObj['totalNumberOfRatings']
        self.likeTotal = mongoObj['likeTotal']
        self.generalRatingTotal = mongoObj['generalRatingTotal']
        self.locationRatingTotal = mongoObj['locationRatingTotal']
        self.eventRatingTotal = mongoObj['eventRatingTotal']
        self.comments = mongoObj['comments']

    def get_event_details_as_json(self):
        user_data = {
            "fullName": self.fullName,
            "age": self.age,
            "gender": self.gender,
            "phoneNumber": self.phoneNumber,
            "email": self.email
        }

        return user_data
