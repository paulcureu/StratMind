package main

import "stratmind/core/parser"

func main() {
	demoPath := "demos/falcons-vs-faze-m3-mirage.dem"
	outputPath := "D:\\VSCode\\StratMind\\core\\export\\timelines.json"

	parser.TrackRounds(demoPath, outputPath)
}

