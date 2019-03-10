import json
from enum import Enum

sampleJson1 = '{"Id":"id1", "message":"Looking for a ride from Brampton to Waterloo on 10th March (Sunday)."}'
sampleJson2 = '{"Id":"id2", "message":"Looking for ride to Union/Finch from Waterloo bk on Sunday (10th) after 4pm."}'
sampleJson3 = '{"Id":"id3", "message":"Driving London -> Waterloo @ 1 pm on Sunday March 10th, $20"}'
sampleJson4 = '{"Id":"id4", "message":"driving richmond hill freshco plaza to waterloo bk plaza at 1pm sunday march 10, no middle seat, taking 407, $20 a seat"}'

locations = {
    ("tdot","toronto","dt") : "Toronto",
    ("rh","richmond hill", "markham") : "Richmond Hill",
    ("loo", "waterloo", "burger king", "bk", "uw") : "Waterloo"
}
def parser(post):
    p = Post(post)
    message = p.message


class Post():
    def __init__(self, jsonObject):
        self.__dict__ = json.loads(jsonObject)

class Trip():
    def __init__(self, pickupLocation, dropoffLocation, pickupTime, driver):
        self.pickupLocation = pickupLocation
        self.dropoffLocation = dropoffLocation
        self.pickupTime = pickupTime
        self.driver = driver

class TripRequest():
    def __init__(self, pickupLocation, dropoffLocation, pickupTime):
        self.pickupLocation = pickupLocation
        self.dropoffLocation = dropoffLocation
        self.pickupTime = pickupTime




parser(sampleJson1)
