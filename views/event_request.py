class AddParticipantRequest:

    def __init__(self, parser):
        self.eventId = checkExistence(parser.eventId)
        self.phoneNumber = checkExistence(parser.phoneNumber)
        self.children = checkExistence(parser.children)
        self.adultInformation = checkExistence(parser.adultInformation)


class RemoveParticipantRequest:

    def __init__(self, parser):
        self.eventId = checkExistence(parser.eventId)

class AddFavouritesRequest:

    def __init__(self, parser):
        self.eventId = checkExistence(parser.eventId)


class RemoveFavouritesRequest:

    def __init__(self, parser):
        self.eventId = checkExistence(parser.eventId)


class GradeEventRequest:

    def __init__(self, parser):
        self.eventId = checkExistence(parser.eventId)
        self.generalRating = checkExistence(parser.generalRating)
        self.locationRating = checkExistence(parser.locationRating)
        self.eventRating = checkExistence(parser.eventRating)
        self.comment = parser.comment


def checkExistence(obj):
    if obj is None: raise TypeError
    return obj
