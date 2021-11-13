import os
import struct
import random
from passlib.hash import pbkdf2_sha256 as sha256


class OTP:
    
    @staticmethod
    def generateRandom():
        #need a crypt safe random generator
        random_number = random.randint(100000, 999999)
        #random_number = struct.unpack('I', os.urandom(length))
        return random_number

    
class Sha:
    """
    Hash password
    """
    @staticmethod
    def generate_hash(password):

        return sha256.hash(password)

    """
    Verify hash and password
    """
    @staticmethod
    def verify_hash(password, hash_):

        return sha256.verify(password, hash_)