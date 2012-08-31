/*
 submit can be used to submit jobs to a mumax-daemon que.
 submit takes a flag describing the kind of job file to make
 and a number of input files. E.g.:
 	submit --mumax2 *.py
 This will generate JSON job files in the user's que directy.
 Default is $HOME/que, can be overridden with --que flag.

 The job files are in simple JSON and may as well be written
 by hand or a script, if more flexibility is needed.

 Author: Arne Vansteenkiste
*/
package main
