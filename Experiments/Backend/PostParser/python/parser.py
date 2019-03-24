import boto3
import json
import re
import csv
import datetime

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

sampleJson1 = '{"Id":"id1", "message":"Looking for ride to Union/Finch from Waterloo bk on Sunday (10th) after 4pm.", "postTime" : "2018-02-31T05:33:31+0000", "username":"Brendan Zhang"}'
sampleJson2 = '{"Id":"id2", "message":"Looking for a ride from Brampton to Waterloo on 10th March (Sunday).", "postTime":"2018-02-31T05:33:31+0000", "username":"Daniell Yang"}'
sampleJson3 = '{"Id":"id3", "message":"Driving London -> Waterloo @ 1 pm on Sunday March 10th, $20", "postTime" : "2018-02-31T05:33:31+0000", "username":"Bimesh DeSilva"}'
sampleJson4 = '{"Id":"id4", "message":"driving richmond hill freshco plaza to waterloo bk plaza at 1pm sunday march 10, no middle seat, taking 407, $20 a seat", "postTime" : "2018-02-31T05:33:31+0000", "username":"Max Gao"}'
shitpost = '{"Id":"id5", "message":"Shitpost", "postTime" : "2018-02-31T05:33:31+0000", "username":"shitposter"}'

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

    pickupLocationSegment = re.search("(?<=DRIVING)(.*)(?=( TO | -> ))", message.upper()).group(1)
    dropoffLocationSegment = message.upper().split(pickupLocationSegment)[1]
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
    print(trip.dropoffLocation)
    print(trip.pickupLocation)
    print(trip.pickupTime)
    print(trip.driver)
    # will actually put this into a db once we have that working

def create_triprequest(message, postTime):
    pickupLocation = ""
    dropoffLocation = ""
    pickupTime = getPickupTime(postTime)

    pickupLocationSegment = re.search("(?<=LOOKING)(.*?)(?=( TO | -> ))", message.upper()).group(1)
    dropoffLocationSegment = message.upper().split(pickupLocationSegment)[1]

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
    print(tripRequest.dropoffLocation)
    print(tripRequest.pickupLocation)
    print(tripRequest.pickupTime)
    #will actually put this into a db once we have that working

def create_manual(message, username, postTime):
    row = [message, username, postTime]
    with open('shitposts.csv','w') as file:
        writer = csv.writer(file)
        writer.writerow(row)
    file.close()

def getPickupTime(postTime):
    year = postTime[0:4]
    month = postTime[5:7]
    day = postTime[8:10]
    #dt = datetime.datetime.strptime(month + " " + day + " " + year, fmt)
    dt = month + " " + day + " " + year
    return dt

def parse_message(message, username, postTime):
    postType = re.search("(DRIV|LOOK)", message.upper())
    if (postType):
        if postType.group(1) == "DRIV":
            create_trip(message, username, postTime)
        if postType.group(1) == "LOOK":
            create_triprequest(message, postTime)
    else:
        create_manual(message, username, postTime)


def parse_json(post):
    p = Post(post)
    parse_message(p.message, p.username, p.postTime)



parse_json(sampleJson2)
#parse_json(sampleJson3)
#parse_json(sampleJson4)
#parse_json(shitpost)
