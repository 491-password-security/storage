from sqlalchemy.dialects.postgresql import ARRAY
from sqlalchemy import *
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
    confirmedAt = sql.Column(sql.DateTime(), nullable=False)
    email = sql.Column(sql.String(320), nullable=True) #max length possible
    active = sql.Column(sql.SmallInteger(), nullable=False)
    password = sql.Column(sql.String(100), nullable=True)
    # User fields
    fullName = sql.Column(sql.String(100), nullable=True)
    profilePictureUrl = sql.Column(sql.String(50), nullable=True)
    username = sql.Column(sql.String(50), nullable=True)
    about = sql.Column(sql.String(320), nullable=True)
    followers = sql.Column(ARRAY(String), nullable=True)
    following = sql.Column(ARRAY(String), nullable=True)
    #social
    network = sql.Column(sql.String(50), nullable=True)
    token = sql.Column(sql.String(320), nullable=True)
    timezone = sql.Column(sql.String(50), nullable=True)

    def get_user_details_as_json(self):
        return {
            "email": self.email,
            "fullName": self.fullName,
            "profilePictureUrl": self.profilePictureUrl,
            "username": self.username,
            "about": self.about,
            "followers": self.followers,
            "following": self.following,
            "network": self.network,
            "token": self.token,
            "timezone": self.timezone
        }

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
    def find_by_email(cls, email):

        return cls.query.filter_by(email=email).first()


    @classmethod
    def find_by_token(cls, token):

        return cls.query.filter_by(token=token).first()

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
