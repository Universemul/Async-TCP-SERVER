package main

const Port = 8001
const Quit = "QUIT"
const Hello = "EHLO"
const Date = "DATE"
const Welcome = "WELCOME"
const DateFormat = "02-01-2006T15:04:05\n"

var availableVerbs = [4]string{Quit, Hello, Date}
