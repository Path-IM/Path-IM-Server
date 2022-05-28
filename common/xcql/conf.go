package xcql

type CassandraConfig struct {
	Hosts         []string
	Port          int
	Keyspace      string
	Username      string
	Password      string
	TimeoutSecond int64
	Consistency   string // Any One Two Three Quorum All LocalQuorum EachQuorum LocalOne
}
