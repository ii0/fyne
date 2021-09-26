// Code generated by go run gen.go; DO NOT EDIT.

package async

import "fyne.io/fyne/v2"

// UnboundedCanvasObjectChan is a channel with an unbounded buffer for caching
// CanvasObject objects.
type UnboundedCanvasObjectChan struct {
	in, out chan fyne.CanvasObject
	close   chan struct{}
}

// NewUnboundedCanvasObjectChan returns a unbounded channel with unlimited capacity.
func NewUnboundedCanvasObjectChan() *UnboundedCanvasObjectChan {
	ch := &UnboundedCanvasObjectChan{
		// The size of CanvasObject is less than 16 bytes, we use 16 to fit
		// a CPU cache line (L2, 256 Bytes), which may reduce cache misses.
		in:    make(chan fyne.CanvasObject, 16),
		out:   make(chan fyne.CanvasObject, 16),
		close: make(chan struct{}),
	}
	go func() {
		// This is a preallocation of the internal unbounded buffer.
		// The size is randomly picked. But if one changes the size, the
		// reallocation size at the subsequent for loop should also be
		// changed too. Furthermore, there is no memory leak since the
		// queue is garbage collected.
		q := make([]fyne.CanvasObject, 0, 1<<10)
		for {
			select {
			case e, ok := <-ch.in:
				if !ok {
					close(ch.out)
					return
				}
				q = append(q, e)
			case <-ch.close:
				goto closed
			}
			for len(q) > 0 {
				select {
				case ch.out <- q[0]:
					q[0] = nil // de-reference earlier to help GC
					q = q[1:]
				case e, ok := <-ch.in:
					if ok {
						q = append(q, e)
						break
					}
					for _, e := range q {
						ch.out <- e
					}
					close(ch.out)
					return
				case <-ch.close:
					goto closed
				}
			}
			// If the remaining capacity is too small, we prefer to
			// reallocate the entire buffer.
			if cap(q) < 1<<5 {
				q = make([]fyne.CanvasObject, 0, 1<<10)
			}
		}

	closed:
		close(ch.in)
		for e := range ch.in {
			q = append(q, e)
		}
		for len(q) > 0 {
			select {
			case ch.out <- q[0]:
				q[0] = nil // de-reference earlier to help GC
				q = q[1:]
			default:
			}
		}
		close(ch.out)
		close(ch.close)
	}()
	return ch
}

// In returns the send channel of the given channel, which can be used to
// send values to the channel.
func (ch *UnboundedCanvasObjectChan) In() chan<- fyne.CanvasObject { return ch.in }

// Out returns the receive channel of the given channel, which can be used
// to receive values from the channel.
func (ch *UnboundedCanvasObjectChan) Out() <-chan fyne.CanvasObject { return ch.out }

// Close closes the channel.
func (ch *UnboundedCanvasObjectChan) Close() {
	ch.close <- struct{}{}
}
