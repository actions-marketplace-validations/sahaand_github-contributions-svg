package main

// neighbors returns all valid adjacent grid positions (up/down/left/right).
func neighbors(p point, numCols, numRows int) []point {
	dirs := []point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	var out []point
	for _, d := range dirs {
		n := point{p.col + d.col, p.row + d.row}
		if n.col >= 0 && n.col < numCols && n.row >= 0 && n.row < numRows {
			out = append(out, n)
		}
	}
	return out
}

// bfsNearest finds the BFS-nearest point among targets from `from`, avoiding blocked positions.
// Returns the target and the full path to it (not including `from`).
func bfsNearest(from point, targets []point, blocked map[point]bool, numCols, numRows int) (point, []point) {
	if len(targets) == 0 {
		return point{}, nil
	}
	targetSet := make(map[point]bool, len(targets))
	for _, t := range targets {
		targetSet[t] = true
	}

	type node struct {
		pos  point
		prev int // index into nodes slice
	}
	visited := make(map[point]bool, 128)
	visited[from] = true
	for b := range blocked {
		visited[b] = true
	}

	var nodes []node
	nodes = append(nodes, node{from, -1})
	queue := []int{0}

	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]
		cur := nodes[idx]
		for _, n := range neighbors(cur.pos, numCols, numRows) {
			if visited[n] {
				continue
			}
			visited[n] = true
			ni := len(nodes)
			nodes = append(nodes, node{n, idx})
			if targetSet[n] {
				var path []point
				for i := ni; i != 0; i = nodes[i].prev {
					path = append([]point{nodes[i].pos}, path...)
				}
				return n, path
			}
			queue = append(queue, ni)
		}
	}
	return point{}, nil
}

// buildHeadPath returns the ordered sequence of positions visited by the snake head.
// The snake uses BFS to eat every non-empty cell, then returns to (0,0).
func buildHeadPath(cells []cell, numCols, numRows int) []point {
	grid := make(map[point]int, len(cells))
	for _, c := range cells {
		if c.level > 0 {
			grid[point{c.col, c.row}] = c.level
		}
	}

	remaining := make(map[point]bool, len(grid))
	for p := range grid {
		remaining[p] = true
	}

	headPath := []point{{0, 0}}

	bodyBlocked := func() map[point]bool {
		k := len(headPath) - 1
		bs := make(map[point]bool, bodyLen)
		for i := 1; i < bodyLen && k-i >= 0; i++ {
			bs[headPath[k-i]] = true
		}
		return bs
	}

	targetSlice := func() []point {
		ts := make([]point, 0, len(remaining))
		for p := range remaining {
			ts = append(ts, p)
		}
		return ts
	}

	origin := point{0, 0}

	for len(remaining) > 0 {
		cur := headPath[len(headPath)-1]
		bs := bodyBlocked()
		targets := targetSlice()

		target, path := bfsNearest(cur, targets, bs, numCols, numRows)
		if path == nil {
			moved := false
			for _, n := range neighbors(cur, numCols, numRows) {
				if !bs[n] {
					headPath = append(headPath, n)
					moved = true
					break
				}
			}
			if !moved {
				break
			}
			continue
		}

		for _, step := range path {
			headPath = append(headPath, step)
			if step == target {
				delete(remaining, target)
				break
			}
		}
	}

	cur := headPath[len(headPath)-1]
	if cur != origin {
		bs := bodyBlocked()
		_, path := bfsNearest(cur, []point{origin}, bs, numCols, numRows)
		if path != nil {
			headPath = append(headPath, path...)
		}
	}

	return headPath
}

// buildEatEvents records which step index the snake head lands on each non-empty cell.
func buildEatEvents(headPath []point, cells []cell) []eatEvent {
	levelMap := make(map[point]int, len(cells))
	for _, c := range cells {
		if c.level > 0 {
			levelMap[point{c.col, c.row}] = c.level
		}
	}
	eaten := make(map[point]bool)
	var events []eatEvent
	for i, p := range headPath {
		if l, ok := levelMap[p]; ok && !eaten[p] {
			events = append(events, eatEvent{i, p, l})
			eaten[p] = true
		}
	}
	return events
}
