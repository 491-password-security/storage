from flask import Flask

from flask_restful import Api

from flask_sqlalchemy import SQLAlchemy

from flask_jwt_extended import JWTManager

from pymongo import MongoClient

# Making Flask Application
app = Flask(__name__)

# Object of Api class
api = Api(app)

# Application Configuration
app.config['SQLALCHEMY_DATABASE_URI'] = 'postgres://postgres:postgres@postgres:5432/jwt_auth'

app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

app.config['SECRET_KEY'] = 'binbuda'

app.config['JWT_SECRET_KEY'] = 'binbesyuzbuda'

app.config['JWT_BLACKLIST_ENABLED'] = True

app.config['JWT_BLACKLIST_TOKEN_CHECKS'] = ['access', 'refresh']

# SqlAlchemy object
sql = SQLAlchemy(app)

# Pymongo cli
#client = MongoClient('mongodb://mongo:mongo@mongo:27017/?authSource=admin')
#db = client.zipzip

#mongo_events = db.events
#mongo_event_themes = db.event_themes

# JwtManager object
jwt = JWTManager(app)


# Generating sql tables before first request is fetched
@app.before_first_request
def create_tables():
    try:
        sql.create_all()
    except:
        print("fatal error creating db")


# Checking that token is in blacklist or not
@jwt.token_in_blacklist_loader
def check_if_token_in_blacklist(decrypted_token):
    jti = decrypted_token['jti']

    return user_models.RevokedTokenModel.is_jti_blacklisted(jti)


# Importing models and resources
from service import auth_service
#from service import event_service
from models import user_models


# Api Endpoints

api.add_resource(auth_service.SimpleEndpointTest, '/test')

api.add_resource(auth_service.UserRegistration, '/registration')

api.add_resource(auth_service.UserLogoutAccess, '/logout/access')

api.add_resource(auth_service.UserLogoutRefresh, '/logout/refresh')

api.add_resource(auth_service.TokenRefresh, '/token/refresh')

api.add_resource(auth_service.PersonalInfo, '/user-info')

api.add_resource(auth_service.SendOtp, '/send-otp')

api.add_resource(auth_service.VerifyOtp, '/verify-otp')