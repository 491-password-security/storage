from collections import UserDict
from enum import unique
import logging
from sqlalchemy import exists
from sqlalchemy.orm import backref
from sqlalchemy.sql.expression import and_, null
from app import sql

from passlib.hash import pbkdf2_sha256 as sha256


class UserModel(sql.Model):
    """
    User Model Class
    """
    __tablename__ = 'user'

    id = sql.Column(sql.Integer, primary_key=True)

    phoneNumber = sql.Column(sql.String(120), unique=True, nullable=False)
    phoneConfirmedAt = sql.Column(sql.DateTime())
    email = sql.Column(sql.String(320)) #max length possible
    active = sql.Column(sql.SmallInteger())
    otp = sql.Column(sql.String(100))
    deviceId = sql.Column(sql.String(50), nullable=False)

    # User fields
    fullName = sql.Column(sql.String(50), nullable=False)
    age = sql.Column(sql.SmallInteger(), nullable=False)
    gender = sql.Column(sql.SmallInteger(), nullable=False)
    permissionToContact = sql.Column(sql.SmallInteger(), nullable=False)
    event_attandee = sql.RelationshipProperty('EventAttandeeModel', backref='user')
    favourited_event = sql.RelationshipProperty('FavouritedEventModel', backref='user')

    def get_user_details_as_json(self):
        
        user_data = {
            "fullName" : self.fullName,
            "age" : self.age,
            "gender" : self.gender,
            "phoneNumber" : self.phoneNumber,
            "email" : self.email
        }

        return user_data

    def update_user(self, data):
        
        for key, value in data.items():
            if(value is not None):
                setattr(self, key, value)

        self.save_to_db()

    """
    Save user details in Database
    """
    def save_to_db(self):

        sql.session.add(self)

        sql.session.commit()

    """
    Commit changes to Database
    """
    def commit_to_db(self):

        sql.session.commit()

    """
    Find user by phone number
    """
    @classmethod
    def find_by_phoneNumber(cls, phoneNumber):

        return cls.query.filter_by(phoneNumber=phoneNumber).first()

    """
    return all the user data in json form available in sql
    """
    @classmethod
    def return_all(cls):

        def to_json(x):

            return {

                'phoneNumber': x.phoneNumber,

            }

        return {'users': [to_json(user) for user in UserModel.query.all()]}

    """
    Delete user data
    """
    @classmethod
    def delete_all(cls):

        try:

            num_rows_deleted = sql.session.query(cls).delete()

            sql.session.commit()

            return {'message': f'{num_rows_deleted} row(s) deleted'}

        except:

            return {'message': 'Something went wrong'}

class RevokedTokenModel(sql.Model):
    """
    Revoked Token Model Class
    """

    __tablename__ = 'revoked_token'

    id = sql.Column(sql.Integer, primary_key=True)

    jti = sql.Column(sql.String(120))

    """
    Save Token in sql
    """
    def add(self):

        sql.session.add(self)

        sql.session.commit()

    """
    Checking that token is blacklisted
    """
    @classmethod
    def is_jti_blacklisted(cls, jti):

        query = cls.query.filter_by(jti=jti).first()

        return bool(query)

class ChildModel(sql.Model):
    
    __tablename__ = 'child'

    id = sql.Column(sql.Integer, primary_key=True)

    parentId = sql.Column(sql.Integer, nullable=False)
    creationDate = sql.Column(sql.DateTime())
    status = sql.Column(sql.SmallInteger())
    
    # User fields
    fullName = sql.Column(sql.String(50), nullable=False)
    age = sql.Column(sql.SmallInteger(), nullable=False)
    gender = sql.Column(sql.SmallInteger(), nullable=False)
    favEvents = sql.Column(sql.ARRAY(sql.Integer))
    attendedEvents = sql.Column(sql.ARRAY(sql.Integer))


    def get_child_details_as_json(self):
        
        user_data = {
            "childId" : self.id,
            "parentId" : self.parentId,
            "fullName" : self.fullName,
            "age" : self.age,
            "gender" : self.gender,
            "status" : self.status,
            "favEvents" : self.favEvents,
            "attendedEvents" : self.attendedEvents
        }

        return user_data

    def update_child(self, data):
        
        for key, value in data.items():
            if(value is not None):
                setattr(self, key, value)

        self.save_to_db()

    """
    Save user details in Database
    """
    def save_to_db(self):

        sql.session.add(self)

        sql.session.commit()

    """
    Commit changes to Database
    """
    def commit_to_db(self):

        sql.session.commit()

    """
    Find user by phone number
    """
    @classmethod
    def find_by_parentId(cls, parentId):

        return cls.query.filter_by(parentId=parentId)

    @classmethod
    def find_by_Id(cls, childId):

        return cls.query.filter_by(id=childId).first()

    """
    return all the user data in json form available in sql
    """
    @classmethod
    def return_all(cls):

        def to_json(x):

            return {

                'phoneNumber': x.phoneNumber,

            }

        return {'users': [to_json(user) for user in UserModel.query.all()]}

    """
    Delete user data
    """
    @classmethod
    def delete_all(cls):

        try:

            num_rows_deleted = sql.session.query(cls).delete()

            sql.session.commit()

            return {'message': f'{num_rows_deleted} row(s) deleted'}

        except:

            return {'message': 'Something went wrong'}

class EventAttandeeModel(sql.Model):
        
    __tablename__ = 'event_attandee'

    id = sql.Column(sql.Integer, primary_key=True, nullable=False, unique=True, autoincrement=True, index=True)
    eventId = sql.Column(sql.String(50), nullable=False)
    parentId = sql.Column(sql.Integer, sql.ForeignKey('user.id'), nullable=False)
    parentAttending = sql.Column(sql.Integer, nullable=False)

    parentName = sql.Column(sql.String(50))
    parentContactPhoneNumber = sql.Column(sql.String(50), nullable=False)
    parentContactEmail = sql.Column(sql.String(50), nullable=False)
    
    childNames = sql.Column(sql.ARRAY(sql.String), nullable=False)
    childAges = sql.Column(sql.ARRAY(sql.Integer), nullable=False)
    childGenders = sql.Column(sql.ARRAY(sql.Integer), nullable=False)

    creationDate = sql.Column(sql.DateTime())
    active = sql.Column(sql.SmallInteger())

    
    def get_event_attendee_details_as_json(self):
        
        user_data = {
            "eventId" : self.eventId,
            "parentId" : self.parentId,
            "parentAttending" : self.parentAttending,
            "active" : self.active,
            "creationDate" : self.creationDate
        }

    def update_event_attandee(self, data):
        
        for key, value in data.items():
            if(value is not None):
                setattr(self, key, value)

        self.save_to_db()

    def save_to_db(self):

        sql.session.add(self)

        sql.session.commit()

    def commit_to_db(self):

        sql.session.commit()

        return cls.query.filter_by(eventId=eventId, parentId=userId).first()

    def delete_from_db(self):

        sql.session.delete(self)

        sql.session.commit()

    @classmethod
    def find_by_eventId(cls, eventId, userId):

        return cls.query.filter_by(eventId=eventId, parentId=userId).first()

class FavouritedEventModel(sql.Model):
        
    __tablename__ = 'favourited_event'

    id = sql.Column(sql.Integer, primary_key=True, nullable=False, unique=True, autoincrement=True, index=True)
    eventId = sql.Column(sql.String(50), nullable=False)
    parentId = sql.Column(sql.Integer, sql.ForeignKey('user.id'), nullable=False)
    creationDate = sql.Column(sql.DateTime())

    def update(self, data):
        
        for key, value in data.items():
            if(value is not None):
                setattr(self, key, value)

        self.save_to_db()

    def save_to_db(self):

        sql.session.add(self)

        sql.session.commit()

    def delete_from_db(self):

        sql.session.delete(self)

        sql.session.commit()

    @classmethod
    def find_by_eventId(cls, eventId, userId):

        return cls.query.filter_by(eventId=eventId, parentId=userId).first()

    @classmethod
    def is_user_favourited_event(cls, userId, eventId):

        favEvent = cls.query.filter_by(eventId=eventId, parentId=userId).first()
        logging.info(favEvent)
        if favEvent is None:
            return False
        else:
            return True