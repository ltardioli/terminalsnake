package main

import "github.com/gdamore/tcell/v2"

const SnakeSymbol = 0x25CF
const AppleSymbol = 0x25CF
const SpecialAppleSymbol = '*'
const GameFrameWidth = 35
const GameFrameHigh = 15
const GameFrameSymbol = '║'
const SpecialAppleChance = 1000 // Chances by frame will be 1/SpecialAppleChance

var screen tcell.Screen
var snake *Snake
var apples []*Apple
var pointsToClear []*Point
var isGamePaused bool
var isGameOver bool
var debugLog string
var score int
var simultaneousApples int
