package slidingbuffer

// note: SlidingBuffer is not thread safe

type SlidingBuffer struct {
	buffer         []byte
	index0         int
	capacity       int
	shrinkMultiple int // at what multiplier of capacity should we shrink buffer?
}

func (sb *SlidingBuffer) Len() int {
	return sb.index0 + len(sb.buffer)
}

func (sb *SlidingBuffer) Append(bytes []byte) {
	sb.buffer = append(sb.buffer, bytes...)
	sb.maybeGc()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// returns the actual start index you're getting, the bytes, and the next index you want
//
// note: we just return a view of our underlying data; this is safe only if
// nobody ever modifies the underlying data
//
// this is slightly redundant cause actualStart + len(data) == nextIndex
func (sb *SlidingBuffer) Get(start int) (int, []byte, int) {
	actualStart := max(start, sb.index0)
	index := actualStart - sb.index0
	data := sb.buffer[index:]
	nextIndex := sb.Len()
	return actualStart, data, nextIndex
}

func (sb *SlidingBuffer) maybeGc() {
	if len(sb.buffer) > sb.capacity*sb.shrinkMultiple {
		// shrink back to capacity
		toRemove := len(sb.buffer) - sb.capacity
		sb.buffer = sb.buffer[toRemove:]
		sb.index0 += toRemove
	}
}

func New(capacity int) *SlidingBuffer {
	sb := &SlidingBuffer{
		buffer:         make([]byte, 0),
		index0:         0,
		capacity:       capacity,
		shrinkMultiple: 2,
	}
	return sb
}
