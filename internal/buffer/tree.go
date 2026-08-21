package buffer

import "unicode/utf8"

const (
	maxLeafPieces    = 64
	maxInternalNodes = 16
)

type treeNode struct {
	parent *treeNode
	leaf   bool

	pieces   []piece
	children []*treeNode

	prev *treeNode
	next *treeNode

	metrics metric
}

type pieceTree struct {
	sources *sources
	root    *treeNode
}

func newPieceTree(sources *sources, pieces []piece) *pieceTree {
	t := &pieceTree{sources: sources}

	if len(pieces) == 0 {
		t.root = &treeNode{leaf: true}
		return t
	}

	leaves := make([]*treeNode, 0, (len(pieces)+maxLeafPieces-1)/maxLeafPieces)

	var prev *treeNode

	for start := 0; start < len(pieces); start += maxLeafPieces {
		end := min(start+maxLeafPieces, len(pieces))

		leaf := &treeNode{
			leaf:   true,
			pieces: append([]piece(nil), pieces[start:end]...),
		}

		leaf.recompute()

		if prev != nil {
			prev.next = leaf
			leaf.prev = prev
		}

		leaves = append(leaves, leaf)
		prev = leaf
	}

	level := leaves

	for len(level) > 1 {
		nextLevel := make(
			[]*treeNode,
			0,
			(len(level)+maxInternalNodes-1)/maxInternalNodes,
		)

		for start := 0; start < len(level); start += maxInternalNodes {
			end := min(start+maxInternalNodes, len(level))

			node := &treeNode{
				children: append([]*treeNode(nil), level[start:end]...),
			}

			for _, child := range node.children {
				child.parent = node
			}

			node.recompute()
			nextLevel = append(nextLevel, node)
		}

		level = nextLevel
	}

	t.root = level[0]

	return t
}

func (n *treeNode) recompute() {
	var total metric

	if n.leaf {
		for _, p := range n.pieces {
			total = addMetric(total, p.metrics)
		}
	} else {
		for _, child := range n.children {
			total = addMetric(total, child.metrics)
		}
	}

	n.metrics = total
}

func (t *pieceTree) totalBytes() int {
	return t.root.metrics.bytes
}

func (t *pieceTree) totalRunes() int {
	return t.root.metrics.runes
}

func (t *pieceTree) totalNewlines() int {
	return t.root.metrics.newlines
}

func (t *pieceTree) locate(offset int) (*treeNode, int, int) {
	offset = min(max(offset, 0), t.totalBytes())

	n := t.root

	for !n.leaf {
		remaining := offset

		for i, child := range n.children {
			if remaining < child.metrics.bytes {
				n = child
				offset = remaining
				break
			}

			remaining -= child.metrics.bytes

			if remaining == 0 {
				if i == len(n.children)-1 {
					n = child
					offset = 0
				} else {
					n = n.children[i+1]
					offset = 0
				}
				break
			}
		}
	}

	remaining := offset

	for i, p := range n.pieces {
		if remaining < p.metrics.bytes {
			return n, i, remaining
		}

		remaining -= p.metrics.bytes

		if remaining == 0 {
			return n, i + 1, 0
		}
	}

	return n, len(n.pieces), 0
}

func (t *pieceTree) splitAt(offset int) (*treeNode, int) {
	leaf, index, inner := t.locate(offset)

	if index == len(leaf.pieces) || inner == 0 {
		return leaf, index
	}

	p := leaf.pieces[index]
	data := t.sources.bytes(p.source, p.start, p.end)

	left := piece{
		source:  p.source,
		start:   p.start,
		end:     p.start + inner,
		metrics: measure(data[:inner]),
	}

	right := piece{
		source:  p.source,
		start:   p.start + inner,
		end:     p.end,
		metrics: measure(data[inner:]),
	}

	pieces := make([]piece, 0, len(leaf.pieces)+1)
	pieces = append(pieces, leaf.pieces[:index]...)
	pieces = append(pieces, left, right)
	pieces = append(pieces, leaf.pieces[index+1:]...)

	leaf.pieces = pieces
	leaf.recompute()
	t.updateMetricsUp(leaf)

	if len(leaf.pieces) > maxLeafPieces {
		t.splitLeaf(leaf)
	}

	leaf, index, _ = t.locate(offset)

	return leaf, index
}

func (t *pieceTree) insertPiece(offset int, p piece) {
	leaf, index := t.splitAt(offset)

	if index > 0 && canMerge(leaf.pieces[index-1], p) {
		leaf.pieces[index-1] = mergePieces(leaf.pieces[index-1], p)
		leaf.recompute()
		t.updateMetricsUp(leaf)
		return
	}

	if index < len(leaf.pieces) && canMerge(p, leaf.pieces[index]) {
		leaf.pieces[index] = mergePieces(p, leaf.pieces[index])
		leaf.recompute()
		t.updateMetricsUp(leaf)
		return
	}

	leaf.pieces = append(leaf.pieces, piece{})
	copy(leaf.pieces[index+1:], leaf.pieces[index:])
	leaf.pieces[index] = p

	leaf.recompute()
	t.updateMetricsUp(leaf)

	if len(leaf.pieces) > maxLeafPieces {
		t.splitLeaf(leaf)
	}
}

func (t *pieceTree) splitLeaf(leaf *treeNode) {
	mid := len(leaf.pieces) / 2

	right := &treeNode{
		leaf:   true,
		pieces: append([]piece(nil), leaf.pieces[mid:]...),
		prev:   leaf,
		next:   leaf.next,
	}

	if leaf.next != nil {
		leaf.next.prev = right
	}

	leaf.next = right
	leaf.pieces = leaf.pieces[:mid]

	leaf.recompute()
	right.recompute()

	parent := leaf.parent

	if parent == nil {
		root := &treeNode{
			children: []*treeNode{leaf, right},
		}

		leaf.parent = root
		right.parent = root

		root.recompute()
		t.root = root

		return
	}

	index := indexOfChild(parent, leaf)

	parent.children = append(parent.children, nil)
	copy(parent.children[index+2:], parent.children[index+1:])
	parent.children[index+1] = right

	right.parent = parent

	parent.recompute()
	t.updateMetricsUp(parent)

	if len(parent.children) > maxInternalNodes {
		t.splitInternal(parent)
	}
}

func (t *pieceTree) splitInternal(node *treeNode) {
	mid := len(node.children) / 2

	right := &treeNode{
		children: append([]*treeNode(nil), node.children[mid:]...),
	}

	for _, child := range right.children {
		child.parent = right
	}

	node.children = node.children[:mid]

	node.recompute()
	right.recompute()

	parent := node.parent

	if parent == nil {
		root := &treeNode{
			children: []*treeNode{node, right},
		}

		node.parent = root
		right.parent = root

		root.recompute()
		t.root = root

		return
	}

	index := indexOfChild(parent, node)

	parent.children = append(parent.children, nil)
	copy(parent.children[index+2:], parent.children[index+1:])
	parent.children[index+1] = right

	right.parent = parent

	parent.recompute()
	t.updateMetricsUp(parent)

	if len(parent.children) > maxInternalNodes {
		t.splitInternal(parent)
	}
}

func (t *pieceTree) deleteRange(start, end int) {
	if start >= end || t.totalBytes() == 0 {
		return
	}

	start = min(max(start, 0), t.totalBytes())
	end = min(max(end, start), t.totalBytes())

	startLeaf, startIndex := t.splitAt(start)
	endLeaf, endIndex := t.splitAt(end)

	if startLeaf == endLeaf {
		startLeaf.pieces = append(
			startLeaf.pieces[:startIndex],
			startLeaf.pieces[endIndex:]...,
		)

		startLeaf.recompute()
		t.updateMetricsUp(startLeaf)

		if len(startLeaf.pieces) == 0 {
			t.removeLeaf(startLeaf)
		}

		return
	}

	startLeaf.pieces = append(
		[]piece(nil),
		startLeaf.pieces[:startIndex]...,
	)

	endLeaf.pieces = append(
		[]piece(nil),
		endLeaf.pieces[endIndex:]...,
	)

	leaves := make([]*treeNode, 0)

	for leaf := startLeaf; leaf != nil; leaf = leaf.next {
		leaves = append(leaves, leaf)

		if leaf == endLeaf {
			break
		}
	}

	for _, leaf := range leaves {
		if leaf == startLeaf || leaf == endLeaf {
			continue
		}

		t.removeLeaf(leaf)
	}

	if len(startLeaf.pieces) == 0 {
		t.removeLeaf(startLeaf)
	} else {
		startLeaf.recompute()
		t.updateMetricsUp(startLeaf)
	}

	if endLeaf != startLeaf && len(endLeaf.pieces) == 0 {
		t.removeLeaf(endLeaf)
	} else if endLeaf != startLeaf {
		endLeaf.recompute()
		t.updateMetricsUp(endLeaf)
	}

	t.coalesceAround(start)
}

func (t *pieceTree) removeLeaf(leaf *treeNode) {
	if leaf.parent == nil {
		leaf.pieces = nil
		leaf.metrics = metric{}
		return
	}

	if leaf.prev != nil {
		leaf.prev.next = leaf.next
	}

	if leaf.next != nil {
		leaf.next.prev = leaf.prev
	}

	parent := leaf.parent
	index := indexOfChild(parent, leaf)

	copy(parent.children[index:], parent.children[index+1:])
	parent.children = parent.children[:len(parent.children)-1]

	if parent == t.root && len(parent.children) == 1 {
		child := parent.children[0]
		child.parent = nil
		t.root = child
		return
	}

	if parent == t.root && len(parent.children) == 0 {
		t.root = &treeNode{leaf: true}
		return
	}

	parent.recompute()
	t.updateMetricsUp(parent.parent)
}

func (t *pieceTree) coalesceAround(offset int) {
	leaf, index, _ := t.locate(offset)

	if index > 0 &&
		index < len(leaf.pieces) &&
		canMerge(leaf.pieces[index-1], leaf.pieces[index]) {

		leaf.pieces[index-1] = mergePieces(
			leaf.pieces[index-1],
			leaf.pieces[index],
		)

		copy(
			leaf.pieces[index:],
			leaf.pieces[index+1:],
		)

		leaf.pieces = leaf.pieces[:len(leaf.pieces)-1]
		leaf.recompute()
		t.updateMetricsUp(leaf)
	}

	if index == 0 &&
		leaf.prev != nil &&
		len(leaf.pieces) > 0 &&
		len(leaf.prev.pieces) > 0 {

		last := len(leaf.prev.pieces) - 1

		if canMerge(leaf.prev.pieces[last], leaf.pieces[0]) {
			leaf.prev.pieces[last] = mergePieces(
				leaf.prev.pieces[last],
				leaf.pieces[0],
			)

			leaf.pieces = leaf.pieces[1:]

			leaf.prev.recompute()
			leaf.recompute()

			t.updateMetricsUp(leaf.prev)

			if len(leaf.pieces) == 0 {
				t.removeLeaf(leaf)
			}
		}
	}
}

func (t *pieceTree) readRange(start, end int, write func([]byte)) {
	if start >= end || t.totalBytes() == 0 {
		return
	}

	start = min(max(start, 0), t.totalBytes())
	end = min(max(end, start), t.totalBytes())

	leaf, index, inner := t.locate(start)
	position := start

	for leaf != nil && position < end {
		if index >= len(leaf.pieces) {
			leaf = leaf.next
			index = 0
			inner = 0
			continue
		}

		p := leaf.pieces[index]
		data := t.sources.bytes(p.source, p.start, p.end)

		take := min(len(data)-inner, end-position)

		write(data[inner : inner+take])

		position += take
		inner = 0
		index++
	}
}

func (t *pieceTree) prefixMetrics(offset int) metric {
	offset = min(max(offset, 0), t.totalBytes())

	if offset == 0 {
		return metric{}
	}

	n := t.root

	var total metric
	remaining := offset

	for !n.leaf {
		for _, child := range n.children {
			if remaining < child.metrics.bytes {
				n = child
				break
			}

			remaining -= child.metrics.bytes
			total = addMetric(total, child.metrics)

			if remaining == 0 {
				return total
			}
		}
	}

	for _, p := range n.pieces {
		if remaining < p.metrics.bytes {
			data := t.sources.bytes(p.source, p.start, p.end)
			partial := measure(data[:remaining])
			return addMetric(total, partial)
		}

		remaining -= p.metrics.bytes
		total = addMetric(total, p.metrics)

		if remaining == 0 {
			return total
		}
	}

	return total
}

func (t *pieceTree) offsetByRunes(target int) int {
	target = min(max(target, 0), t.totalRunes())

	if target == 0 {
		return 0
	}

	if target == t.totalRunes() {
		return t.totalBytes()
	}

	n := t.root
	offset := 0
	remaining := target

	for !n.leaf {
		for _, child := range n.children {
			if remaining < child.metrics.runes {
				n = child
				break
			}

			remaining -= child.metrics.runes
			offset += child.metrics.bytes

			if remaining == 0 {
				return offset
			}
		}
	}

	for _, p := range n.pieces {
		if remaining < p.metrics.runes {
			data := t.sources.bytes(p.source, p.start, p.end)
			byteOffset := byteOffsetAfterRunes(data, remaining)
			return offset + byteOffset
		}

		remaining -= p.metrics.runes
		offset += p.metrics.bytes

		if remaining == 0 {
			return offset
		}
	}

	return offset
}

func (t *pieceTree) offsetByNewlines(target int) int {
	target = min(max(target, 0), t.totalNewlines())

	if target == 0 {
		return 0
	}

	n := t.root
	offset := 0
	remaining := target

	for !n.leaf {
		for _, child := range n.children {
			if remaining <= child.metrics.newlines {
				n = child
				break
			}

			remaining -= child.metrics.newlines
			offset += child.metrics.bytes
		}
	}

	for _, p := range n.pieces {
		if remaining > p.metrics.newlines {
			remaining -= p.metrics.newlines
			offset += p.metrics.bytes
			continue
		}

		data := t.sources.bytes(p.source, p.start, p.end)

		for i, b := range data {
			if b == '\n' {
				remaining--

				if remaining == 0 {
					return offset + i + 1
				}
			}
		}

		offset += len(data)
	}

	return offset
}

func (t *pieceTree) updateMetricsUp(node *treeNode) {
	for node != nil {
		node.recompute()
		node = node.parent
	}
}

func indexOfChild(parent, child *treeNode) int {
	for i, current := range parent.children {
		if current == child {
			return i
		}
	}

	panic("buffer: child not found")
}

func byteOffsetAfterRunes(data []byte, count int) int {
	if count <= 0 {
		return 0
	}

	offset := 0

	for count > 0 && offset < len(data) {
		_, size := utf8.DecodeRune(data[offset:])
		offset += size
		count--
	}

	return offset
}
