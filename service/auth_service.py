from random import expovariate
from re import T
import time

from flask_restful import Resource, reqparse

from models.user_models import UserModel, RevokedTokenModel
from flask_jwt_extended import (
    create_access_token,
    create_refresh_token,
    jwt_required,
    jwt_refresh_token_required,
    get_jwt_identity,
    get_raw_jwt
)

import pdb
from utility.validation import Validate
from utility.crypt import OTP
from utility.crypt import Sha
from utility.email_ops import Mail

import datetime

import logging

# provide simple and uniform access to any variable
from views.general_response import SuccessResponse, FailureResponse, AuthErrorResponse

parser = reqparse.RequestParser()
parser.add_argument('password')
parser.add_argument('email')
parser.add_argument('fullName')
parser.add_argument('network')
parser.add_argument('token')
parser.add_argument('timezone')
parser.add_argument('profilePicture')
parser.add_argument('username')


class SimpleEndpointTest(Resource):
    def get(self):
        return SuccessResponse("Test Success").__dict__


class SimpleEndpointTestWithAuth(Resource):
    def get(self):
        return SuccessResponse("Test Success").__dict__


class UserRegistration(Resource):
    def post(self):

        data = parser.parse_args()

        logging.info("Post request to '/registration'.")

        email = data['email'].lower();

        if not email:
            logging.error("email in the request body is null.")
            return FailureResponse("email in the request appears to be Null.").__dict__

        # Checking that user is already exist or not
        elif UserModel.find_by_email(email):
            logging.error(f'{email} is already registered.')
            return FailureResponse(f'{email} is already registered').__dict__


        try:
            password = OTP.generateRandom()

            Mail.sendEmail(email, password)

        except:
            logging.error("something went wrong")
            return FailureResponse('Something went wrong').__dict__


        new_user = UserModel(

            active=1,
            email=email,
            fullName=data['fullName'],
            password=str(Sha.generate_hash(str(password))),
            confirmedAt=datetime.datetime.now()
        )


        try:
            new_user.save_to_db()

            logging.info("User info saved to db.")
            return SuccessResponse({}).__dict__

        except:
            logging.info(
                "Error occured saving user info to db, might be caused of some recured fields being empty in the request body.")
            return FailureResponse("Error occured saving user info to db").__dict__


class VerifyOtp(Resource):

    def post(self):

        data = parser.parse_args()

        logging.info("Post request to '/login'.")

        email = data['email'].lower();

        if not email:
            logging.error("email in the request body is null.")
            return FailureResponse('phoneNumber in the request appears to be Null.').__dict__

        # Searching user by email
        current_user = UserModel.find_by_email(email)

        # user does not exists
        if not current_user:
            logging.error(f'User with email {email} doesn\'t exist.')
            return FailureResponse(f'User with email {email} doesn\'t exist').__dict__

        logging.error(str(Sha.generate_hash(data['password'])))
        logging.error(current_user.password)

        if Sha.verify_hash(data['password'], current_user.password):

            try:

                access_token = create_access_token(identity=email)

                refresh_token = create_refresh_token(identity=email)

                return SuccessResponse({

                    'access_token': access_token,

                    'refresh_token': refresh_token

                }).__dict__

            except:

                return FailureResponse('Something went wrong').__dict__

        return FailureResponse('Password or Username does not match').__dict__


class UserLogoutAccess(Resource):
    """
    User Logout Api 
    """

    @jwt_required
    def post(self):

        jti = get_raw_jwt()['jti']

        try:
            # Revoking access token
            revoked_token = RevokedTokenModel(jti=jti)

            revoked_token.add()

            return SuccessResponse("Access token has been revoked").__dict__

        except:

            return FailureResponse("Failed to logout").__dict__


class UserLogoutRefresh(Resource):
    """
    User Logout Refresh Api 
    """

    @jwt_refresh_token_required
    def post(self):

        jti = get_raw_jwt()['jti']

        try:

            revoked_token = RevokedTokenModel(jti=jti)

            revoked_token.add()

            pdb.set_trace()

            return SuccessResponse("Refresh token has been revoked").__dict__

        except:

            return FailureResponse("Something went wrong").__dict__


class TokenRefresh(Resource):
    """
    Token Refresh Api
    """

    @jwt_refresh_token_required
    def get(self):
        # Generating new access token
        try:
            current_user = get_jwt_identity()
        except:
            return FailureResponse("Could not find an account related to this refresh token").__dict__
        
        try:
            access_token = create_access_token(identity=current_user)
        except:
            return FailureResponse("Could not create access token with the given refresh token").__dict__


        return SuccessResponse({'access_token': access_token}).__dict__


class PersonalInfo(Resource):
    @jwt_required
    def get(self):

        try:
            data = parser.parse_args()
        except:
            return FailureResponse("Error parsing args").__dict__

        try:
            email = get_jwt_identity()
        except:
            return AuthErrorResponse("Access Token invalid").__dict__
        # Searching user by phoneNumber
        
        try:
            current_user = UserModel.find_by_email(email)
        except:
            return FailureResponse("Can not find the current users email in db").__dict__

        # user does not exists
        if not current_user:
            return FailureResponse(f'User with email {email} doesn\'t exist').__dict__
        else:
            return SuccessResponse(current_user.get_user_details_as_json()).__dict__

    @jwt_required
    def post(self):

        try:
            data = parser.parse_args()
        except:
            return FailureResponse("Error parsing args").__dict__

        try:
            email = get_jwt_identity()
        except:
            return AuthErrorResponse("Access Token invalid").__dict__


        try:
            current_user = UserModel.find_by_email(email)
        except:
            return FailureResponse("Can not find the current users email in db").__dict__

        if not current_user:
            return {'message': f'User with email {email} doesn\'t exist'}

        try:
            current_user.update_user(data)
        except:
            return FailureResponse("Can not update the user").__dict__

        return SuccessResponse({}).__dict__


class SocialLogin(Resource):

    def post(self):
        try:
            data = parser.parse_args()
        except:
            return FailureResponse("Error parsing args").__dict__

        token = data['token']
        try:
            current_user = UserModel.find_by_token(token)
        except:
            return FailureResponse("Can not find the current users token in db").__dict__

        if not current_user:

            new_user = UserModel(
                active=1,
                confirmedAt=datetime.datetime.now(),
                token=token,
                timezone=data['timezone'],
                network=data['network'],
                email=token
            )

            try:
                new_user.save_to_db()

                access_token = create_access_token(identity=token)

                refresh_token = create_refresh_token(identity=token)

                return SuccessResponse({

                    'access_token': access_token,

                    'refresh_token': refresh_token

                }).__dict__

            except:
                return FailureResponse('Something went wrong').__dict__
        try:

            access_token = create_access_token(identity=token)

            refresh_token = create_refresh_token(identity=token)

            return SuccessResponse({

                'access_token': access_token,

                'refresh_token': refresh_token

            }).__dict__

        except:
            return FailureResponse('Something went wrong').__dict__


class SendOtp(Resource):
    def post(self):

        update_dict = {}

        data = parser.parse_args()

        email = data['email']

        current_user = UserModel.find_by_email(email)

        if not current_user:
            return FailureResponse("There is no account connected to this email address").__dict__
        try:
            password = OTP.generateRandom()

            Mail.sendEmail(email, password)

            update_dict['password'] = str(Sha.generate_hash(str(password)))

            current_user.update_user(update_dict)
            logging.info("updated the user")
        except:
            logging.error("something went wrong")
            return FailureResponse('Something went wrong').__dict__

        return SuccessResponse('Success').__dict__
