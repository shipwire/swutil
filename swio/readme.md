# swio
--
    import "github.com/shipwire/swutil/swio"

Package swio provides additional utilities on top of the standard io package.

## Usage

```go
var DummyReader = newDummy(time.Now().Unix())
```
DummyReader is a reader of pseudo-random data. It is meant to be more efficient
than cryptographically random data, but is useful only in limited cases such as
testing other readers.

#### func  ForkReader

```go
func ForkReader(r io.Reader, n int) (head, tail io.Reader)
```
ForkReader accepts a reader and forks it at a given point. It returns two
readers, one that continues from the beginning of the original reader, and
another that reads from the nth byte. Readers are a FIFO stream, so reads from
the second reader can't actually jump ahead. To solve this problem, reads from
the tail cause the entire remaining contents of the head to be transparently
read into memory.

#### func  NewReadCloser

```go
func NewReadCloser(r io.Reader, c CloseFunc) io.ReadCloser
```
NewReadCloser wraps r with CloseFunc c.

#### func  SeekerPosition

```go
func SeekerPosition(r io.Seeker) (int64, error)
```
SekerPosition returns the current offset of an io.Seeker

#### func  TeeReadCloser

```go
func TeeReadCloser(r io.ReadCloser, w io.Writer) io.ReadCloser
```
TeeReadCloser returns a ReadCloser that writes everything read from it to w. Its
content is read from r. All writes must return before anything can be read. Any
write error will be returned from Read. The Close method on the returned
ReadCloser is called derived from r.

#### type CloseFunc

```go
type CloseFunc func() error
```

CloseFunc documents a function that satisfies io.Closer.

#### type ReadCounter

```go
type ReadCounter interface {
	io.Reader
	BytesRead() int64
}
```

ReadCounter is a reader that keeps track of the total number of bytes read.

#### func  NewReadCounter

```go
func NewReadCounter(r io.Reader) ReadCounter
```
NewReadCounter wraps a Reader in a ReadCounter.

#### type ReadSeekerAt

```go
type ReadSeekerAt interface {
	io.ReadSeeker
	io.ReaderAt
}
```

ReadSeekerAt is the interface that groups io.ReadSeeker and io.ReaderAt.

#### type ReadSeekerCloser

```go
type ReadSeekerCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
```

ReadSeekerCloser wraps the io.Reader, io.Seeker, and io.Closer types.

#### func  NewReadSeekerCloser

```go
func NewReadSeekerCloser(r io.ReadSeeker, c CloseFunc) ReadSeekerCloser
```
NewReadSeekerCloser wraps r with the close function c.

#### func  NewSeekBufferCloser

```go
func NewSeekBufferCloser(r io.ReadCloser, l int64) ReadSeekerCloser
```
NewSeekBufferCloser creates a new SeekBuffer using ReadCloser r. l is the
initial capacity of the internal buffer. The close method from r is forwarded to
the returned ReadSeekerCloser.

#### func  NopReadSeekerCloser

```go
func NopReadSeekerCloser(r io.ReadSeeker) ReadSeekerCloser
```
NopReadSeekerCloser wraps r with a no-op close function.

#### type SeekBuffer

```go
type SeekBuffer struct {
}
```

SeekBuffer contains a reader where all data read from it is buffered, and thus
both readable and seekable. It essentially augments bytes.Reader so that all
data does not need to be read at once.

#### func  NewSeekBuffer

```go
func NewSeekBuffer(r io.Reader, l int64) *SeekBuffer
```
NewSeekBuffer creates a new SeekBuffer using reader r. l is the initial capacity
of the internal buffer. If r happens to be a ReadSeeker already, r is used
directly without any additional copies of the data.

#### func (*SeekBuffer) Read

```go
func (s *SeekBuffer) Read(b []byte) (int, error)
```
Read reads from the internal buffer or source reader until b is filled or an
error occurs.

#### func (*SeekBuffer) ReadAt

```go
func (s *SeekBuffer) ReadAt(b []byte, off int64) (int, error)
```
ReadAt implements io.ReaderAt for SeekBuffer.

#### func (*SeekBuffer) Seek

```go
func (s *SeekBuffer) Seek(offset int64, whence int) (int64, error)
```
Seek implements io.Seeker for SeekBuffer. If the new offset goes beyond what is
buffered, SeekBuffer will read from the source reader until offset can be
reached or an error is returned. If whence is 2, SeekBuffer will read until EOF
or an error is reached.

#### type WriteCounter

```go
type WriteCounter interface {
	io.Writer
	BytesWritten() int64
}
```

WriteCounter is a writer that keeps track of the total number of bytes written.

#### func  NewWriteCounter

```go
func NewWriteCounter(w io.Writer) WriteCounter
```
NewWriteCounter wraps a Writer in a WriteCounter
