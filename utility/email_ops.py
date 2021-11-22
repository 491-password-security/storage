import smtplib, ssl
import logging

import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from mailjet_rest import Client

class Mail:
    @staticmethod
    def sendEmail(to, new_pass):
        api_key = '4135fbd0efbcf74c63f2d70c956ccc8a'
        api_secret = '70fc15e4740bba0d771593dde77b622f'
        mailjet = Client(auth=(api_key, api_secret), version='v3.1')
        logging.error(new_pass)
        data = {
          'Messages': [
            {
              "From": {
                "Email": "info.veviski@gmail.com",
                "Name": "Veviski"
              },
              "To": [
                {
                  "Email": to,
                  "Name": ""
                }
              ],
              "Subject": "Veviski Password Reset",
              "TextPart": str(new_pass),
              "HTMLPart": "<h3>Your password has been reset</h3>" + str(new_pass),
              "CustomID": "AppGettingStartedTest"
            }
          ]
        }
        result = mailjet.send.create(data=data)
        logging.error(result.status_code)
        logging.error(result.json())


