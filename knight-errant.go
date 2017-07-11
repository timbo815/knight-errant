package main

import(
  "html/template"
  "net/http"
  "encoding/json"
)

const (
  boardWidth = 8
  boardHeight = 8
  numSquares = 64
)

type ResponseObject struct {
	Response [boardWidth][boardHeight]int `json:"response"`
}

func renderTemplate(w http.ResponseWriter, tmpl string, varMap interface{}) {
  t, _ := template.ParseFiles(tmpl + ".html")
  t.Execute(w, varMap)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  varMap := map[string]interface{}{
    "answer": "",
  }
  renderTemplate(w, "root", varMap)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
  var board [boardWidth][boardHeight]int
  var visitedSquareCount int
  moveKnight(w, board, visitedSquareCount, 3, 4)
}

func moveKnight(w http.ResponseWriter, board [boardWidth][boardHeight]int, visitedSquareCount int, x int, y int) bool {
  if x >= boardWidth || y >= boardHeight || x < 0 || y < 0 {
    return false
  }

  if board[x][y] > 0 {
    return false
  }

  visitedSquareCount ++
  board[x][y] = visitedSquareCount

  if visitedSquareCount == numSquares {
    ro := &ResponseObject{}
    ro.Response = board

    json, err := json.Marshal(ro)
    if err != nil {
      // log error
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(json)

    return true
  }

  var deltas = [8][2]int{
    {2, 1},
    {2, -1},
    {-2, 1},
    {-2, -1},
    {1, 2},
    {1, -2},
    {-1, 2},
    {-1, -2},
  }

  for _, delta := range deltas {
    if moveKnight(w, board, visitedSquareCount, x + delta[0], y + delta[1]) {
      return true
    }
  }

  board[x][y] = 0
  visitedSquareCount --
  return false
}

func main() {
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/calculate/", calculateHandler)
  http.ListenAndServe(":8080", nil)
}
