import datetime
from http.client import PROXY_AUTHENTICATION_REQUIRED
import pprint

import pymongo

from flask_restful import Resource, reqparse
from pymongo.message import query
from app import mongo_events, mongo_event_themes
from utility.navigation import getCity, calculateDistance
from views import event_response, event_request
from bson import json_util, ObjectId
import json
import logging
from queue import PriorityQueue
from models.user_models import FavouritedEventModel, UserModel, RevokedTokenModel, ChildModel, EventAttandeeModel


from views.general_response import SuccessResponse, FailureResponse

logging.basicConfig(level=logging.DEBUG)


class SimpleMongoTest(Resource):
    def get(self):
        event = mongo_events.find_one()

        return SuccessResponse(parse_json(event)).__dict__

class GetThemes(Resource):
    def get(self):
        themes = list(mongo_event_themes.find({}))
        return SuccessResponse(parse_json(themes)).__dict__

class GetEventDetailsById(Resource):
    def get(self):
        from flask import request
        phoneNumber = request.headers.get('phone')

        parser = reqparse.RequestParser()
        parser.add_argument('eventId', type=str, required=True)

        eventId = parser.parse_args().eventId

        event = mongo_events.find_one({'_id': ObjectId(eventId)})
        current_user = UserModel.find_by_phoneNumber(phoneNumber)

        event['favourited'] = FavouritedEventModel.is_user_favourited_event(current_user.id, eventId)

        response = event_response.EventDetailResponseModel(event)

        return SuccessResponse(response.__dict__).__dict__

class GetEventList(Resource):
    def get(self):

        parser = reqparse.RequestParser()
        parser.add_argument('city')
        parser.add_argument('type')
        parser.add_argument('themes')
        parser.add_argument('gender')
        parser.add_argument('age')
        parser.add_argument('distance')
        parser.add_argument('sortBy')
        parser.add_argument('sortDirection')

        args = parser.parse_args()

        city = checkIfEmpty(args.get('city'))
        type = checkIfEmpty(args.get('type'))
        themes = checkIfEmpty(args.get('themes'))
        gender = checkIfEmpty(args.get('gender'))
        age = checkIfEmpty(args.get('age'))
        distance = checkIfEmpty(args.get('distance'))
        sortBy = (args.get("sortBy"))
        sortDirection = (args.get("sortDirection"))

        query = {"city": city, "type": type, "themes": themes, "gender": gender}

        logging.debug(query)

        if sortBy is None:
            events = mongo_events.find(query)
        else:

            if sortDirection is None:
                return FailureResponse("While sorting you have to specify direction").__dict__

            if str(sortDirection) == "ascending":
                events = mongo_events.aggregate([{
                    "$match": query
                },
                    {
                        "$sort": {
                            sortBy: pymongo.ASCENDING,
                        }
                    }])
            elif str(sortDirection) == "descending":
                events = mongo_events.aggregate([{
                    "$match": query
                },
                    {
                        "$sort": {
                            sortBy: pymongo.DESCENDING,
                        }
                    }])
            else:
                return FailureResponse("Wrong format for sortDirection in request body").__dict__

        if events is None:
            logging.error("events is none")

        response = []
        for event in events:
            result = event_response.EventBasicResponseModel(event)
            response.append(result.__dict__)

        return SuccessResponse(parse_json(response)).__dict__

class GetEventListWithDistance(Resource):
    def get(self):

        parser = reqparse.RequestParser()
        parser.add_argument('city')
        parser.add_argument('type')
        parser.add_argument('themes')
        parser.add_argument('gender')
        parser.add_argument('age')
        parser.add_argument('distance')
        parser.add_argument('sortBy')
        parser.add_argument('sortDirection')
        parser.add_argument('coordinates')

        args = parser.parse_args()

        city = (args.get('city'))
        type = checkIfEmpty(args.get('type'))
        themes = checkIfEmpty(args.get('themes'))
        gender = checkIfEmpty(args.get('gender'))
        age = checkIfEmpty(args.get('age'))
        distance = checkIfEmpty(args.get('distance'))
        sortBy = args.get("sortBy")
        sortDirection = args.get("sortDirection")
        clientCoordinates = args.get("coordinates")
        clientCoordinateX, clientCoordinateY = clientCoordinates.split(",")

        if city is None:
            city = getCity(float(clientCoordinateX), float(clientCoordinateY))
            logging.info("client is in " + str(city))

        query = {"city": city, "type": type, "themes": themes, "gender": gender}

        logging.debug(query)

        if sortBy is None:
            events = mongo_events.find(query)
        else:

            if sortDirection is None:
                return FailureResponse("While sorting you have to specify direction").__dict__

            if str(sortDirection) == "ascending":
                events = mongo_events.aggregate([{
                    "$match": query
                },
                    {
                        "$sort": {
                            sortBy: pymongo.ASCENDING,
                        }
                    }])
            elif str(sortDirection) == "descending":
                events = mongo_events.aggregate([{
                    "$match": query
                },
                    {
                        "$sort": {
                            sortBy: pymongo.DESCENDING,
                        }
                    }])
            else:
                return FailureResponse("Wrong format for sortDirection in request body").__dict__

        if events is None:
            logging.error("events is none")

        response = []
        for event in events:
            result = event_response.EventWithDistanceResponseModel(event)
            logging.debug(result.__dict__)
            eventCoordinateX, eventCoordinateY = event['coordinates'].split(",")
            distance = calculateDistance(float(clientCoordinateX), float(clientCoordinateY),
                                         float(eventCoordinateX), float(eventCoordinateY))
            result.distance = distance
            response.append(result.__dict__)

        response.sort(key=lambda x: x['distance'])

        return SuccessResponse(parse_json(response)).__dict__

class AddParticipant(Resource):
    def post(self):
        from flask import request
        parser = reqparse.RequestParser()
        # parser.add_argument('eventId')

        phoneNumber = request.headers.get('phone')
        # parser.add_argument('children') # 1 or many 
        # parser.add_argument('adultInformation')

        json_data = request.get_json(force=True)
        logging.info(json_data)
        eventId = json_data.get('eventId')
        children = json_data.get('children')
        adultInformation = json_data.get('adultInformation')
        if eventId is None or children is None or adultInformation is None:
            return FailureResponse("Missing parameters in request body").__dict__
        
        logging.info(phoneNumber)
        current_user = UserModel.find_by_phoneNumber(phoneNumber)
        if current_user is None:
            return FailureResponse("User not found").__dict__
        

        childNames = []
        childAges = []
        childGenders = []
        
        for child in children:
            if child is None:
                return FailureResponse("Missing child in request body").__dict__
            childNames.append(child.get('name'))
            childAges.append(int(child.get('age')))
            childGenders.append(int(child.get('gender')))

        newTicket = EventAttandeeModel(

            eventId = eventId,
            parentId = int(current_user.id),

            parentAttending = int(adultInformation.get('attending')),
            parentName = adultInformation.get('name'),
            parentContactEmail = adultInformation.get('contactEmail'),
            parentContactPhoneNumber = adultInformation.get('contactPhoneNumber'),

            childNames = childNames,
            childAges = childAges,
            childGenders = childGenders,

            creationDate = datetime.datetime.now(),
            active = 1
        )

        try:
            newTicket.save_to_db()
        except:
            return FailureResponse("Error in adding participant").__dict__

        return SuccessResponse("Participant added to event " + eventId).__dict__

class RemoveParticipant(Resource):
    def post(self):
        from flask import request
        phoneNumber = request.headers.get('phone')

        parser = reqparse.RequestParser()
        parser.add_argument('eventId')

        try:
            request = event_request.RemoveParticipantRequest(parser.parse_args())
        except:
            return FailureResponse("Body does not fit the required structure").__dict__

        user = UserModel.find_by_phoneNumber(phoneNumber)

        EventAttandeeModel.find_by_eventId(request.eventId, user.id).delete_from_db()

        return SuccessResponse("Participant removed from event " + request.eventId).__dict__

class GetRegisteredEventList(Resource):
    def get(self):
        from flask import request
        phoneNumber = request.headers.get('phone')

        user = UserModel.find_by_phoneNumber(phoneNumber)

        attendedEvents = user.event_attandee
        #logging.info(attendedEvents[0].__dict__)
        eventIds = []
        for ticket in attendedEvents:
            eventIds.append(ObjectId(ticket.eventId))

        logging.info(eventIds)
        
        query = {"_id": {"$in": eventIds}}
        events = list(mongo_events.aggregate([{
                    "$match": query
                }]))

        # for event in events:
        #     logging.info(event.__dict__)

        logging.info(events)

        response = []
        for event in events:
            result = event_response.EventBasicResponseModel(event)
            response.append(result.__dict__)
        # return SuccessResponse(parse_json(response)).__dict__
        return SuccessResponse(parse_json(response)).__dict__

class GetFavouriteEvents(Resource):
    def get(self):
        from flask import request
        phoneNumber = request.headers.get('phone')

        user = UserModel.find_by_phoneNumber(phoneNumber)

        favouritedEvents = user.favourited_event

        eventIds = []
        for favEvents in favouritedEvents:
            eventIds.append(ObjectId(favEvents.eventId))

        logging.info(eventIds)
        
        query = {"_id": {"$in": eventIds}}
        events = list(mongo_events.aggregate([{
                    "$match": query
                }]))

        logging.info(events)

        response = []
        for event in events:
            result = event_response.EventBasicResponseModel(event)
            response.append(result.__dict__)
        # return SuccessResponse(parse_json(response)).__dict__
        return SuccessResponse(parse_json(response)).__dict__

class AddFavourites(Resource):
    def post(self):
        from flask import request
        phoneNumber = request.headers.get('phone')

        json_data = request.get_json(force=True)
        eventId = json_data.get('eventId')

        if eventId is None:
            return FailureResponse("Missing parameters in request body").__dict__
        
        current_user = UserModel.find_by_phoneNumber(phoneNumber)
        if current_user is None:
            return FailureResponse("User not found").__dict__

        newFavEvent = FavouritedEventModel(
            eventId = eventId,
            parentId = int(current_user.id),
            creationDate = datetime.datetime.now()
        )

        try:
            newFavEvent.save_to_db()
        except:
            return FailureResponse("Error in adding participant").__dict__

        return SuccessResponse("Favourited event:  " + eventId).__dict__

class RemoveFavourites(Resource):
    def post(self):
        from flask import request
        phoneNumber = request.headers.get('phone')

        parser = reqparse.RequestParser()
        parser.add_argument('eventId')

        try:
            request = event_request.RemoveParticipantRequest(parser.parse_args())
        except:
            return FailureResponse("Body does not fit the required structure").__dict__

        user = UserModel.find_by_phoneNumber(phoneNumber)

        FavouritedEventModel.find_by_eventId(request.eventId, user.id).delete_from_db()

        return SuccessResponse("Unfavourited event: " + request.eventId).__dict__

class GradeEvent(Resource):
    def post(self):
        parser = reqparse.RequestParser()
        parser.add_argument('eventId')
        parser.add_argument('generalRating')
        parser.add_argument('locationRating')
        parser.add_argument('eventRating')
        parser.add_argument('comment')

        try:
            request = event_request.GradeEventRequest(parser.parse_args())
        except:
            return FailureResponse("Body does not fit the required structure").__dict__

        return SuccessResponse("Participant graded the event " + request.eventId).__dict__

def checkIfEmpty(arg):
    if arg is None:
        return {'$regex': '.*'}
    return arg

def parse_json(data):
    return json.loads(json_util.dumps(data))
