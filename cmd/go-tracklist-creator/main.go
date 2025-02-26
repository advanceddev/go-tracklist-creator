package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Pair - запись типа Артист - Название
type Pair struct {
	Artist string
	Track  string
}

// Graph - граф для построения пути
type Graph map[string][]string

func main() {

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	file, err := os.Open("drops.txt")
	if err != nil {
		logger.Fatalf("Не удалось открыть drops.txt: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pairs []Pair

	var currentPair Pair
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if currentPair.Artist != "" && currentPair.Track != "" {
				pairs = append(pairs, currentPair)
				currentPair = Pair{}
			}
			continue
		}

		parts := strings.SplitN(line, " - ", 2)
		if len(parts) != 2 {
			logger.Warnf("Неверный формат записи: %s", line)
			continue
		}

		if currentPair.Artist == "" {
			currentPair.Artist = parts[0]
			currentPair.Track = parts[1]
		} else {
			pairs = append(pairs, currentPair)
			currentPair.Artist = parts[0]
			currentPair.Track = parts[1]
		}
	}

	if currentPair.Artist != "" && currentPair.Track != "" {
		pairs = append(pairs, currentPair)
	}

	if err := scanner.Err(); err != nil {
		logger.Fatalf("Не удалось прочитать файл drops.txt: %v", err)
	}

	if len(pairs) == 0 {
		logger.Fatalf("В файле drops.txt не найдено подходящих пар треков для генерации треклиста.")
	}

	graph := CreateGraph(pairs)
	ShuffleGraph(graph)

	tracklist := FindTracklist(graph, pairs[0].Artist+" - "+pairs[0].Track)

	err = WriteTracklistToFile(tracklist, "tracklist.txt")
	if err != nil {
		logger.Fatalf("Не удалось сохранить tracklist.txt: %v", err)
	}

	logger.Println("Треклист сгенерирован и записан в файл tracklist.txt")
}

// CreateGraph - создает граф на основе пар треков
func CreateGraph(pairs []Pair) Graph {
	graph := make(Graph)
	seen := make(map[string]bool)

	for _, pair := range pairs {
		key := pair.Artist + " - " + pair.Track
		if seen[key] {
			continue
		}
		seen[key] = true
		if graph[key] == nil {
			graph[key] = []string{}
		}
	}

	for i := 0; i < len(pairs)-1; i++ {
		currentKey := pairs[i].Artist + " - " + pairs[i].Track
		nextKey := pairs[i+1].Artist + " - " + pairs[i+1].Track
		if graph[currentKey] == nil {
			graph[currentKey] = []string{}
		}
		graph[currentKey] = append(graph[currentKey], nextKey)
	}

	if len(pairs) > 0 {
		lastKey := pairs[len(pairs)-1].Artist + " - " + pairs[len(pairs)-1].Track
		if !seen[lastKey] {
			graph[lastKey] = []string{}
		}
	}

	return graph
}

// ShuffleGraph - случайно перемешивает соседние треки в графе
func ShuffleGraph(graph Graph) {
	rand.Seed(time.Now().UnixNano())

	for node, neighbors := range graph {
		if len(neighbors) > 1 {
			rand.Shuffle(len(neighbors), func(i, j int) {
				neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
			})
			graph[node] = neighbors
		}
	}
}

// FindTracklist - находит путь в графе без повторений
func FindTracklist(graph Graph, start string) []string {
	var tracklist []string
	visited := make(map[string]bool)

	var dfs func(node string)
	dfs = func(node string) {
		if visited[node] {
			return
		}
		visited[node] = true
		tracklist = append(tracklist, node)
		for _, neighbor := range graph[node] {
			dfs(neighbor)
		}
	}

	dfs(start)
	return tracklist
}

// WriteTracklistToFile - записывает треклист в файл
func WriteTracklistToFile(tracklist []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, track := range tracklist {
		_, err := writer.WriteString(track + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
