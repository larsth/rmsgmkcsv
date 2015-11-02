#!/bin/python3

from sys import stdin
from sys import stdout
import codecs

def rmsgmkcsvio():
	#Read from standard input (STDIN) ...
	degreesStr = stdin.readline()
	#	input_stream = codecs.getreader("utf-8")(stdin)
	#	degreesStr = input_stream.readline()
	#Translate the 'degreesStr' string to the 'degreesF' 
	#floating-point value:
	degreesF = float(degreesStr)
	#Calculate a microcontroller value using the 'degreesF' 
	#floating-point value ...
	resultF = degreesF * 123
	#Translate the 'resultF' floating-point microcontroller 
	#value to an integer ...
	resultI = int(resultF)
	#Translate the 'resultI' integer microcontroller value to a string ...
	resultS = str(resultI)
	resultS += "\n" #renember to add a trailing '\n' (the split token)
	#Write the 'resultS' string to standard output (STDOUT) ...
	stdout.write(resultS)
	stdout.flush()

#Below: only used for invoke once for testing
#uncommented when not used ...
#rmsgmkcsvio()

#loop until killed by the parent process ...
while True:
	rmsgmkcsvio()
	
	
