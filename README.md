# Gator
A RSS aggregator that will let you browse headlines and descriptions from the command line. 

## Use
To use this program you will require Postgres and Go to be installed on your computer. After running go build you can run the program from your command line with ./gator

## Setup
All functions need to be run with ./gator, this will be omited for the setup

You will need to register your account with 
register [name]
This will be the name that feeds are saved under

login [name]
This will switch the feed view from one user to another

reset
This will clear all feeds for everyone

users
This will provide you with a list of all users on the local machine

agg [[number][timeFrame]] eg: agg 60s
This will poll all followed RSS feeds every [number][timeFrame]. In the example above it would be very 60 seconds

addfeed [name] [url]
Adds a feed under name at url

feeds
Prints out all feeds you are following

unfollow [url]
Removes url from yoru following list

browse [?number]
Takes an optional number that will return number or 2 articles that were most recently updated