package main

import "fmt"

const DownScore = 1
const LeftScore = 2
const RightScore = 0
const UpScore = 3
const DEPTH = 10

var TOTAL int = 0
var RESULTS map[string]map[string]int = map[string]map[string]int{}

type Card struct {
	Name  string
	Score [4]int
}

type Player []*Card

type Score map[string]int

func (p Player) HasCards() bool {
	return len(p) > 0
}

type Players map[string]Player

type Position struct {
	Owner string
	Card  *Card
}

type Board [3][3]Position

func (b *Board) IsFull() bool {
	for _, rows := range b {
		for _, cell := range rows {
			if cell.Card == nil {
				return false
			}
		}
	}

	return true
}

func main() {
	p := Players{
		"P1": Player{
			&Card{"C0", [4]int{10, 3, 9, 5}},
			//&Card{"C1", [4]int{2, 3, 5, 7}},
			&Card{"C2", [4]int{4, 6, 10, 7}},
			//&Card{"C3", [4]int{5, 5, 6, 4}},
			//&Card{"C4", [4]int{4, 10, 4, 8}},
		},
		"P2": Player{
			//&Card{"C5", [4]int{3, 6, 7, 2}},
			&Card{"C6", [4]int{5, 4, 5, 6}},
			//&Card{"C7", [4]int{8, 8, 4, 2}},
			//&Card{"C8", [4]int{1, 4, 3, 6}},
			//&Card{"C9", [4]int{1, 1, 2, 6}},
		},
	}

	b := Board{
		[3]Position{
			Position{
				"P2",
				&Card{"C7", [4]int{8, 8, 4, 2}},
			},
			Position{
				"P1",
				&Card{"C8", [4]int{1, 4, 3, 6}},
			},
			Position{
				"P1",
				&Card{"C3", [4]int{6, 6, 7, 5}},
			},
		},
		[3]Position{
			Position{},
			Position{
				"P1",
				&Card{"C1", [4]int{1, 2, 4, 6}},
			},
			Position{
				"P1",
				&Card{"C4", [4]int{4, 10, 4, 8}},
			},
		},
		[3]Position{
			Position{},
			Position{
				"P2",
				&Card{"C9", [4]int{1, 1, 2, 6}},
			},
			Position{
				"P2",
				&Card{"C5", [4]int{3, 6, 7, 2}},
			},
		},
	}

	calculate(b, p, "P1")
	for kres, vres := range RESULTS {
		p1 := vres["P1"]
		p2 := vres["P2"]
		draw := vres["draw"]
		total := vres["total"]
		fmt.Printf(`
%s
-----------
P1: %d (%.2f)
P2: %d (%.2f)
Draw: %d (%.2f)
TOTAL: %d

		`, kres,
			p1, 100*float64(p1)/float64(total),
			p2, 100*float64(p2)/float64(total),
			draw, 100*float64(draw)/float64(total),
			total)
	}

	fmt.Printf("TOTAL: %d\n", TOTAL)
	printBoard(b)

}
func calculate(b Board, p Players, n string) {
	calculateWithDepth(b, p, n, DEPTH, "")
}

func calculateWithDepth(b Board, p Players, n string, depth int, strategy string) {
	TOTAL = TOTAL + 1
	//printBoard(b)
	//fmt.Println(p)
	//calculateScore(b, p)
	//fmt.Println(depth)
	//fmt.Println(n)
	nextPlayer := "P1"
	if n == "P1" {
		nextPlayer = "P2"
	}
	if b.IsFull() || !p[n].HasCards() || depth <= 0 {
		score := calculateScore(b, p)
		updateResult(score, strategy)
		return
	}

	for nc, c := range p[n] {
		for idy, x := range b {
			for idx, pos := range x {
				if pos.Card != nil {
					continue
				}
				tmpBoard := b
				tmpPlayers := Players{}
				for pk, pv := range p {
					tmpPlayers[pk] = Player{}
					for _, cv := range pv {
						tmpPlayers[pk] = append(tmpPlayers[pk], cv)
					}
				}
				tmpBoard[idy][idx].Card = c
				tmpBoard[idy][idx].Owner = n
				tmpPlayers[n][nc] = tmpPlayers[n][len(tmpPlayers[n])-1]
				tmpPlayers[n] = tmpPlayers[n][:len(tmpPlayers[n])-1]
				tmpBoard = resolveBoard(tmpBoard, idx, idy)
				newStrategy := strategy
				if strategy == "" {
					newStrategy = fmt.Sprintf("%s-%d-%d", c.Name, idx, idy)
					RESULTS[newStrategy] = map[string]int{}
				}
				calculateWithDepth(tmpBoard, tmpPlayers, nextPlayer, depth-1, newStrategy)
			}
		}
	}
}

func resolveBoard(b Board, x int, y int) Board {
	if x < 2 && b[y][x+1].Card != nil && b[y][x+1].Card.Score[LeftScore] < b[y][x].Card.Score[RightScore] {
		b[y][x+1].Owner = b[y][x].Owner
	}

	if x > 0 && b[y][x-1].Card != nil && b[y][x-1].Card.Score[RightScore] < b[y][x].Card.Score[LeftScore] {
		b[y][x-1].Owner = b[y][x].Owner
	}

	if y < 2 && b[y+1][x].Card != nil && b[y+1][x].Card.Score[UpScore] < b[y][x].Card.Score[DownScore] {
		b[y+1][x].Owner = b[y][x].Owner
	}

	if y > 0 && b[y-1][x].Card != nil && b[y-1][x].Card.Score[DownScore] < b[y][x].Card.Score[UpScore] {
		b[y-1][x].Owner = b[y][x].Owner
	}

	return b
}

func calculateScore(b Board, ps Players) Score {
	res := map[string]int{}
	for pn, p := range ps {
		s := len(p)
		for _, rows := range b {
			for _, cell := range rows {
				if cell.Owner == pn {
					s = s + 1
				}
			}
		}
		res[pn] = s
	}

	return res
}

func updateResult(score Score, strategy string) {
	if _, ok := RESULTS[strategy]["P1"]; !ok {
		RESULTS[strategy]["P1"] = 0
	}
	if _, ok := RESULTS[strategy]["P2"]; !ok {
		RESULTS[strategy]["P2"] = 0
	}
	if _, ok := RESULTS[strategy]["total"]; !ok {
		RESULTS[strategy]["total"] = 0
	}
	if _, ok := RESULTS[strategy]["draw"]; !ok {
		RESULTS[strategy]["draw"] = 0
	}
	if score["P1"] > score["P2"] {
		RESULTS[strategy]["P1"] += 1
	}
	if score["P1"] < score["P2"] {
		RESULTS[strategy]["P2"] += 1
	}
	if score["P1"] == score["P2"] {
		RESULTS[strategy]["draw"] += 1
	}
	RESULTS[strategy]["total"] += 1
}

func printBoard(b Board) {
	v := []interface{}{}
	for _, rows := range b {
		for _, cell := range rows {
			if cell.Card == nil {
				v = append(v, 0)
				continue
			}
			v = append(v, cell.Card.Score[UpScore])
		}
		for _, cell := range rows {
			if cell.Card == nil {
				v = append(v, 0)
				v = append(v, "  ")
				v = append(v, 0)
				continue
			}
			v = append(v, cell.Card.Score[LeftScore])
			v = append(v, cell.Owner)
			v = append(v, cell.Card.Score[RightScore])
		}
		for _, cell := range rows {
			if cell.Card == nil {
				v = append(v, 0)
				continue
			}
			v = append(v, cell.Card.Score[DownScore])
		}
	}
	fmt.Printf(
		`
-------------------------
|  %d    |  %d    |  %d    |
|%d %s %d |%d %s %d |%d %s %d |
|  %d    |  %d    |  %d    |
------------------------- 
|  %d    |  %d    |  %d    |
|%d %s %d |%d %s %d |%d %s %d |
|  %d    |  %d    |  %d    |
------------------------- 
|  %d    |  %d    |  %d    |
|%d %s %d |%d %s %d |%d %s %d |
|  %d    |  %d    |  %d    |
-------------------------
		`, v...,
	)
}
