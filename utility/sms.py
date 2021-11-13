import logging

import requests
import sys


def sendOtp(otp):
    response = requests.get(
        'https://api.netgsm.com.tr/sms/send/otp',
        params={'usercode': '03246060458', 'password': 'ZipZip21', 'no': '5345001428',
                'msg': 'Şifreniz: ' + str(otp),
                'msgheader': 'EFG TURIZM'},
    )

    # with open('filename.txt', 'w') as f:
    #     sys.stdout = f  # Change the standard output to the file we created.
    #     print(response)
    logging.DEBUG(response)
    return response.text


class Sms:
    pass
