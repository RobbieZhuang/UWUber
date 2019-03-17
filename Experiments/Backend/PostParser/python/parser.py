import boto3
import json
import re
import csv

dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table('TripPosts')

#response = table.query()

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

sampleJson1 = '{"Id":"id1", "message":"Looking for a ride from Brampton to Waterloo on 10th March (Sunday).", "postTime":"2018-02-31T05:33:31+0000", "username":"Daniell Yang"}'
sampleJson2 = '{"Id":"id2", "message":"Looking for ride to Union/Finch from Waterloo bk on Sunday (10th) after 4pm.", "postTime" : "2018-02-31T05:33:31+0000", "username":"Brendan Zhang"}'
sampleJson3 = '{"Id":"id3", "message":"Driving London -> Waterloo @ 1 pm on Sunday March 10th, $20", "postTime" : "2018-02-31T05:33:31+0000", "username":"Bimesh DeSilva"}'
sampleJson4 = '{"Id":"id4", "message":"driving richmond hill freshco plaza to waterloo bk plaza at 1pm sunday march 10, no middle seat, taking 407, $20 a seat", "postTime" : "2018-02-31T05:33:31+0000", "username":"Max Gao")'

locations = {
    ("TDOT","TORONTO") : "Toronto",
    ("RH","RICHMOND HILL", "MARKHAM") : "Richmond Hill",
    ("LOO", "WATERLOO", "BURGER KING PLAZA", "BK PLAZA" "BURGER KING", "UW", "KW") : "Waterloo",
    ("STC", "SCARBOROUGH") : "Scarborough",
    ("BRAMPTON") : "Brampton",
    ("LONDON", "WESTERN") : "London"
}

def create_trip(message, username, postTime):
    driver = username
    pickupLocation = ""
    dropoffLocation = ""
    pickupTime = getPickupTime(postTime)

    pickupLocationSegment = re.search("(?<=DRIVING).*?(?=(TO|->))", message.upper())
    dropoffLocationSegment = re.search("(?<=(TO|->)).*?(?=(ON|AT|@))", message.upper())
    for tuple in locations.keys():
        for name in tuple:
            if name in pickupLocationSegment:
                pickupLocation = locations[tuple]
                break
            if name in dropoffLocationSegment:
                dropoffLocation = locations[tuple]
                break

    if not pickupLocation or not dropoffLocation:
        create_manual(message, username, postTime)

    trip = Trip(pickupLocation, dropoffLocation, pickupTime, driver)

def create_triprequest(message, postTime):
    pickupLocation = ""
    dropoffLocation = ""
    pickupTime = getPickupTime(postTime)

    pickupLocationSegment = re.search("(?<=LOOKING).*?(?=(TO|->))", message.upper())
    dropoffLocationSegment = re.search("(?<=(TO|->)).*?", message.upper())
    for tuple in locations.keys():
        for name in tuple:
            if name in pickupLocationSegment:
                pickupLocation = locations[tuple]
                break
            if name in dropoffLocationSegment:
                dropoffLocation = locations[tuple]
                break

    if not pickupLocation or not dropoffLocation:
        create_manual(message, "", postTime)

    tripRequest = TripRequest(pickupLocation, dropoffLocation, pickupTime)

def create_manual(message, username, postTime):
    row = [message, username, postTime]
    with open('shitposts.csv','w') as file:
        writer = csv.writer(file)
        writer.writerow(row)
    file.close()

def getPickupTime(postTime):
    #This part is kinda hard so it just returns for now
    return postTime

def parse_message(message, username, postTime):
    postType = re.search("DRIV|LOOK", message.upper())
    if postType == "DRIV":
        create_trip(message, username, postTime)
    elif postType == "LOOK":
        create_triprequest(message, postTime)
    else:
        create_manual(message, username, postTime)


def parse_json(post):
    p = Post(post)
    parse_message(p.message, p.username, p.postTime)



#parse_json(sampleJson1)
