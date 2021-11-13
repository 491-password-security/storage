import re

class Validate:
    
    def isTurkishPhoneNumber(phoneNumber):
        return re.match(r'^(05)\d{9}$', phoneNumber)
