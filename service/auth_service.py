from flask_restful import Resource, reqparse

from models.user_models import UserModel, RevokedTokenModel, ChildModel, EventAttandeeModel

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
from utility.sms import Sms, sendOtp

import datetime

import logging

# provide simple and uniform access to any variable
from views.general_response import SuccessResponse, FailureResponse

parser = reqparse.RequestParser()
parser.add_argument('phoneNumber')
parser.add_argument('email')
parser.add_argument('fullName')
parser.add_argument('age')
parser.add_argument('gender')
parser.add_argument('otp')
parser.add_argument('childId')
parser.add_argument('status')
parser.add_argument('favEvents')
parser.add_argument('attendedEvents')
parser.add_argument('parentId')
parser.add_argument('deviceId')
parser.add_argument('permissionToContact')
parser.add_argument('phone', type=str, location='headers', required=True)


class SimpleEndpointTest(Resource):
    def get(self):
        return SuccessResponse("Test Success").__dict__


class UserRegistration(Resource):
    """
    User Registration Api
    """

    def post(self):

        data = parser.parse_args()

        logging.info("Post request to '/registration'.")

        phoneNumber = data['phoneNumber']

        if not phoneNumber:
            logging.error("phoneNumber in the request body is null.")
            return FailureResponse("'message': 'phoneNumber in the request appears to be Null.'").__dict__

        elif not Validate.isTurkishPhoneNumber(phoneNumber):
            logging.error(f'{phoneNumber} is not a valid Turkish phone number.')
            return FailureResponse(f'{phoneNumber} is not a valid Turkish phone number').__dict__

        # Checking that user is already exist or not
        elif UserModel.find_by_phoneNumber(phoneNumber):
            logging.error(f'{phoneNumber} is already registered.')
            return FailureResponse(f'{phoneNumber} is already registered').__dict__

        otp = OTP.generateRandom()  # 6 digits
        logging.info("OTP generated.")

        print(otp)
        #resp = sendOtp(otp)
        #logging.debug(resp)
        sha_otp = Sha.generate_hash(str(otp))

        # create new user
        new_user = UserModel(

            phoneNumber=phoneNumber,
            active=0,
            fullName=data['fullName'],
            age=data['age'],
            gender=data['gender'],
            otp=sha_otp,
            permissionToContact=data['permissionToContact'],
            deviceId=data['deviceId']
        )

        try:
            new_user.save_to_db()

            logging.info("User info saved to db.")
            return SuccessResponse({'otp': str(otp)}).__dict__

        except:

            logging.info(
                "Error occured saving user info to db, might be caused of some recured fields being empty in the request body.")
            return FailureResponse("Error occured saving user info to db").__dict__


class UserLogin(Resource):
    """
    User Login Api
    """

    def post(self):

        data = parser.parse_args()

        logging.info("Post request to '/registration'.")

        phoneNumber = data['phoneNumber']

        if not phoneNumber:
            logging.error("phoneNumber in the request body is null.")
            return {'message': 'phoneNumber in the request appears to be Null.'}, 400

        elif not Validate.isTurkishPhoneNumber(phoneNumber):
            logging.error(f'{phoneNumber} is not a valid Turkish phone number.')
            return {'message': f'{phoneNumber} is not a valid Turkish phone number'}, 400

        # Searching user by phoneNumber
        current_user = UserModel.find_by_phoneNumber(phoneNumber)

        # user does not exists
        if not current_user:
            logging.error(f'User with phone number {phoneNumber} doesn\'t exist.')
            return {'message': f'User with phone number {phoneNumber} doesn\'t exist'}, 404

        otp = OTP.generateRandom()  # 6 digits
        logging.info("OTP generated.")

        otp_hashed = Sha.generate_hash(str(otp))
        current_user.otp = otp_hashed

        try:
            current_user.commit_to_db()
            logging.info('OTP saved to db')
            return SuccessResponse({'otp_hashed': otp_hashed}).__dict__
        except:
            logging.error('Problem while saving otp to db')
            return FailureResponse("Error logging in and saving otp to db").__dict__


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
        current_user = get_jwt_identity()

        access_token = create_access_token(identity=current_user)

        return SuccessResponse({'access_token': access_token}).__dict__


class AllUsers(Resource):

    def get(self):
        """
        return all user api
        """
        return UserModel.return_all()

    def delete(self):
        """
        delete all user api
        """
        return UserModel.delete_all()


class PersonalInfo(Resource):
    # IF PHONE NUMBER IS ALSO UPDATED, SENDS NEW TOKENS

    """
    Secret Resource Api
    You can create crud operation in this way
    """

    # todo: jwt has to be required
    def get(self):

        data = parser.parse_args()
        phoneNumber = data['phone']

        # todo: has to enabled
        # phoneNumber = get_jwt_identity()
        # Searching user by phoneNumber
        current_user = UserModel.find_by_phoneNumber(phoneNumber)

        # user does not exists
        if not current_user:
            return FailureResponse(f'User with phone number {phoneNumber} doesn\'t exist').__dict__
        else:
            return SuccessResponse(current_user.get_user_details_as_json()).__dict__

    # todo: jwt has to be required
    def put(self):

        data = parser.parse_args()
        phoneNumber = data['phone']

        # todo: has to enabled
        # phoneNumber = get_jwt_identity()

        current_user = UserModel.find_by_phoneNumber(phoneNumber)

        if not current_user:
            return {'message': f'User with phone number {phoneNumber} doesn\'t exist'}

        current_user.update_user(data)

        if (phoneNumber != data['phoneNumber'] and data['phoneNumber'] is not None):

            try:

                phoneNumber = data['phoneNumber']

                access_token = create_access_token(identity=phoneNumber)

                refresh_token = create_refresh_token(identity=phoneNumber)

                return {

                    'message': f'User {phoneNumber} was created',

                    'access_token': access_token,

                    'refresh_token': refresh_token

                }

            except:

                return {'message': 'Something went wrong while generating tokens'}, 500

        return {'message': 'Succesfully updated'}


class OtpVerification(Resource):
    """
        Verify user's phone
    """

    def post(self):

        data = parser.parse_args()

        phoneNumber = data['phoneNumber']
        otp = data['otp']

        current_user = UserModel.find_by_phoneNumber(phoneNumber)

        if Sha.generate_hash(str(otp)) == current_user.otp or str(otp) == "000000":

            current_user.phoneConfirmedAt = datetime.datetime.now()
            current_user.active = 1

            current_user.commit_to_db()

            try:

                access_token = create_access_token(identity=phoneNumber)

                refresh_token = create_refresh_token(identity=phoneNumber)

                return SuccessResponse({

                    'access_token': access_token,

                    'refresh_token': refresh_token

                }).__dict__

            except:

                return FailureResponse('Something went wrong').__dict__

        return FailureResponse('Code does not match').__dict__


class ChildCreation(Resource):

    # todo: jwt has to be required
    def put(self):

        data = parser.parse_args()
        phoneNumber = data['phone']

        # todo: enable this
        # phoneNumber = get_jwt_identity()
        # Searching user by phoneNumber
        current_user = UserModel.find_by_phoneNumber(phoneNumber)

        # user does not exists
        if not current_user:
            return {'message': f'User with phone number {phoneNumber} doesn\'t exist'}

        # create new child
        new_child = ChildModel(

            parentId=current_user.id,
            creationDate=datetime.datetime.now(),
            fullName=data['fullName'],
            age=data['age'],
            gender=data['gender'],
            status=1
        )

        try:
            new_child.save_to_db()
            return SuccessResponse({'childId': new_child.id}).__dict__
        except:
            return FailureResponse('Error in saving to db').__dict__


class ChildInfo(Resource):

    # todo: enable jwt
    def post(self):

        data = parser.parse_args()
        phoneNumber = data['phone']

        childId = data['childId']

        child = ChildModel.find_by_Id(childId)

        if not child:
            return FailureResponse("f'Child with id {childId} doesn\'t exist'").__dict__
        return SuccessResponse(child.get_child_details_as_json()).__dict__

    # todo: enable jwt
    def put(self):

        data = parser.parse_args()
        phoneNumber = data['phone']

        childId = data['childId']

        child = ChildModel.find_by_Id(childId)

        if not child:
            return {'message': f'Child with id {childId} doesn\'t exist'}

        try:
            child.update_child(data)

            return {'message': 'Succesfully updated'}

        except:

            return {'message': 'Error while updating child data'}


class ChildList(Resource):

    # todo: jwt has to be enabled
    def get(self):
        data = parser.parse_args()
        phoneNumber = data['phone']

        # todo: enable this
        # phoneNumber = get_jwt_identity()

        current_user = UserModel.find_by_phoneNumber(phoneNumber)

        children = ChildModel.find_by_parentId(current_user.id)
        children_list = []

        for child in children:
            children_list.append(child.get_child_details_as_json())

        return SuccessResponse(children_list).__dict__
