# swio
--
    import "github.com/shipwire/swutil/swio"


## Usage

```go
var Dummy = NewDummy(time.Now().Unix())
```
Dummy is a reader of pseudo-random data. It is meant to be faster than
cryptographically random data, but is useful in limited cases such as testing
other readers.

#### func  ForkReader

```go
func ForkReader(r io.Reader, n int) (io.Reader, io.Reader)
```
ForkReader accepts a reader and forks it at a given point. It returns two
readers, one that continues from the beginning of the original reader, and
another that reads from the nth byte. Reads from the second reader cause the
entire remaining contents of the first reader to be read into memory.

#### func  NewDummy

```go
func NewDummy(seed int64) io.Reader
```
NewDummy returns a new reader of pseudo-random data. Its returned reader is
meant to be faster than cryptographically random data, but is useful in limited
cases such as testing other readers. Use the same seed value to produce
deterministic data.

#### func  NewReadCloser

```go
func NewReadCloser(r io.Reader, c CloseFunc) io.ReadCloser
```

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
# httpio
--
    import "bitbucket.org/shipwire/swio/httpio"

Package httpio contains io utilities that are particularly relevant to http.

## Usage

```go
const (
	Invalid = iota
	Informational
	Success
	Redirect
	ClientError
	ServerError
)
```

```go
var Discard http.ResponseWriter = discardResponse{}
```
Discard is an http.ResponseWriter on which all Write calls succeed without doing
anything and no headers are stored.

```go
var UTF8Server = new(utf8server)
```
UTF8Server is an http.Handler which replaces a request body with one translated
to UTF-8 if a different content type is specified. If the content type is
unsupported, missing, or malformed, the request will be rejected. See
http://godoc.org/code.google.com/p/go-charset/charset for more information about
supported character sets.

#### func  AggregateStatus

```go
func AggregateStatus(ws ...TrackedResponseWriter) int
```
AggregateStatus attempts to return a reasonable response code that generalizes
an arbitrary set of response codes.

#### func  CombineHeaders

```go
func CombineHeaders(hs ...http.Header) http.Header
```
CombineHeaders accepts a list of headers and combines them into a single map.
Later specified arguments take precedence over earlier ones.

#### type Client

```go
type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
	Head(url string) (*http.Response, error)
	Post(url string, bodyType string, body io.Reader) (*http.Response, error)
	PostForm(url string, data url.Values) (*http.Response, error)
}
```

Client encapsulates an interface the mirrors http.Client.

#### type ResponseBuffer

```go
type ResponseBuffer struct {
}
```

ResponseBuffer is an in-memory implementation of http.ResponseWriter, similar to
bytes.Buffer.

#### func  NewResponseBuffer

```go
func NewResponseBuffer() *ResponseBuffer
```
NewResponseBuffer returns a new ResponseBuffer.

#### func (ResponseBuffer) AddHeadersTo

```go
func (r ResponseBuffer) AddHeadersTo(w http.ResponseWriter)
```
AddHeadersTo copies the responce buffer's headers to another response writer.

#### func (ResponseBuffer) ContentType

```go
func (r ResponseBuffer) ContentType() string
```
ContentType attempts to determine the content type of the response. It first
checks the Content-Type header, and if empty, it uses up to the first 512 bytes
of the response body and checks that with http.DetectContentType.

#### func (*ResponseBuffer) Header

```go
func (r *ResponseBuffer) Header() http.Header
```
Header returns the http.Header map for the response.

#### func (*ResponseBuffer) Status

```go
func (r *ResponseBuffer) Status() int
```
Status returns the current status code for the response.

#### func (*ResponseBuffer) Write

```go
func (r *ResponseBuffer) Write(b []byte) (int, error)
```
Write adds data to the buffer, growing it as needed.

#### func (*ResponseBuffer) WriteHeader

```go
func (r *ResponseBuffer) WriteHeader(i int)
```
Write header sets the status code to the response. Only the first call to this
function does anything.

#### func (*ResponseBuffer) WriteTo

```go
func (r *ResponseBuffer) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements http.WriterTo

#### func (*ResponseBuffer) WriteToResponseWriter

```go
func (r *ResponseBuffer) WriteToResponseWriter(w http.ResponseWriter) (int64, error)
```
WriteToResponseWriter writes the response buffer to another http.ResponseWriter
by first copying headers, then the response code, and then the response body.

#### func (*ResponseBuffer) Written

```go
func (r *ResponseBuffer) Written() bool
```
Written returns true if the response buffer has been written to.

#### type ResponseLevel

```go
type ResponseLevel int
```

ResponseLevel is the response code level, e.g. 2xx, 3xx.

#### func  Level

```go
func Level(code int) ResponseLevel
```
Level takes a response code and returns the ResponseLevel.

#### type StatusAggregate

```go
type StatusAggregate int
```

StatusAggregate is an http status code that may be aggregated with others.

#### func (*StatusAggregate) Add

```go
func (s *StatusAggregate) Add(other int)
```
Add aggregates the current status code with an additional status code, changing
the original, if necessary, to reflect the multiple statuses.

#### type TrackedResponseWriter

```go
type TrackedResponseWriter interface {
	http.ResponseWriter
	Status() int
	Written() bool
}
```

TrackedResponseWriter implements http.ResponseWriter while tracking if data has
been written and which Status Code has been written.

#### func  NewTrackedResponseWriter

```go
func NewTrackedResponseWriter(w http.ResponseWriter) TrackedResponseWriter
```
NewTrackedResponseWriter returns a new TrackedResponseWriter.
# swio
--
    import "bitbucket.org/shipwire/swutil/swio"

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
