package toany

type Client struct {
	S  int
	S1 *string
	S2 []*string
	S3 [][]*string
	S4 [][]*string

	M1 map[string]string
	M2 map[int]string
}
