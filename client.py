#!/usr/bin/env python3

import urllib.request
import json
import time

url = 'https://scrape.asg.cx/api/scrapes'
auth_header = 'Bearer scrape-YWb36X7FvLWVK93crQec9neaJji3M3Ud+50JAtdTNpo'

# Each time you call the API you will likely get scrape files you have already
# processed so you need to keep track of the files you have seen so you don't
# process them a second time. Each scrape file has a unique key so you can
# keep a list of keys you have already seen and ignore any files whose key is
# in the list.
seen = []

if __name__ == '__main__':
	# We want to query the API continuously to get new scrape files.
	while True:
		# Each request to the API must include an authorization header with
		# your bearer token.
		headers = {'Authorization': auth_header}
		res = urllib.request.urlopen(urllib.request.Request(url, headers=headers))

		# The response will be a JSON document that contains a list of
		# scrape files, each of which look like the following:
		#
		# {
		#    'key': 'Z9Ji/f/2Uy/mT38KOHX/mRp5W+9Nfl+4rSks8uqeCQE',
		#    'location': '',
		#    'size': 131,
		#    'user': 'yulei',
		#    'domain': 'github.com',
		#    'sha256': '34dcfdb4d756318cb04ba7ca9b6cb4872ae7a447a403df9c56eb00e93576be94',
		#    'tags': ['amazon-aws-secret', 'match-regex']
		# }
		#
		# You can use the metadata from the scrape file to determine which
		# interesting files you want to download.
		files = json.loads(res.read().decode('utf-8'))

		# Loop through our scrape files to find interesting content
		for file in files:
			# Skip scrape files we've already seen
			if file['key'] in seen:
				continue

			# Print the location of any scrape file that is tagged 'emails'
			if 'emails' in file['tags']:
				print(file['location'])

			# Append the key to our seen list so we don't process this file
			# again.
			seen.append(file['key'])

		# There is no need to query the API more often than every 30 seconds.
		print('Sleeping for 30 seconds...')
		time.sleep(30)
