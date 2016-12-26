package quantizer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-trace-agent/model"
)

type redisTestCase struct {
	query            string
	expectedResource string
}

func RedisSpan(query string) model.Span {
	return model.Span{
		Resource: query,
		Type:     "redis",
	}
}

func TestRedisQuantizer(t *testing.T) {
	assert := assert.New(t)

	queryToExpected := []redisTestCase{
		{"CLIENT LIST",
			"CLIENT LIST"},

		{"get my_key",
			"GET"},

		{"SET le_key le_value",
			"SET"},

		{"\n\n  \nSET foo bar  \n  \n\n  ",
			"SET"},

		{"CONFIG SET parameter value",
			"CONFIG SET"},

		{"SET toto tata \n \n  EXPIRE toto 15  ",
			"PIPELINE [ EXPIRE SET ]"},

		{"MSET toto tata toto tata toto tata \n ",
			"MSET"},

		{"MULTI\nSET k1 v1\nSET k2 v2\nSET k3 v3\nSET k4 v4\nDEL to_del\nEXEC",
			"PIPELINE [ DEL EXEC MULTI SET ]"},

		{"DEL k1\nDEL k2\nHMSET k1 \"a\" 1 \"b\" 2 \"c\" 3\nHMSET k2 \"d\" \"4\" \"e\" \"4\"\nDEL k3\nHMSET k3 \"f\" \"5\"\nDEL k1\nDEL k2\nHMSET k1 \"a\" 1 \"b\" 2 \"c\" 3\nHMSET k2 \"d\" \"4\" \"e\" \"4\"\nDEL k3\nHMSET k3 \"f\" \"5\"\nDEL k1\nDEL k2\nHMSET k1 \"a\" 1 \"b\" 2 \"c\" 3\nHMSET k2 \"d\" \"4\" \"e\" \"4\"\nDEL k3\nHMSET k3 \"f\" \"5\"\nDEL k1\nDEL k2\nHMSET k1 \"a\" 1 \"b\" 2 \"c\" 3\nHMSET k2 \"d\" \"4\" \"e\" \"4\"\nDEL k3\nHMSET k3 \"f\" \"5\"",
			"PIPELINE [ DEL HMSET ]"},
	}

	for _, testCase := range queryToExpected {
		assert.Equal(testCase.expectedResource, Quantize(RedisSpan(testCase.query)).Resource)
	}

}

func BenchmarkTestRedisQuantizer(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		span := RedisSpan(`DEL k1\nDEL k2\nHMSET k1 "a" 1 "b" 2 "c" 3\nHMSET k2 "d" "4" "e" "4"\nDEL k3\nHMSET k3 "f" "5"\nDEL k1\nDEL k2\nHMSET k1 "a" 1 "b" 2 "c" 3\nHMSET k2 "d" "4" "e" "4"\nDEL k3\nHMSET k3 "f" "5"\nDEL k1\nDEL k2\nHMSET k1 "a" 1 "b" 2 "c" 3\nHMSET k2 "d" "4" "e" "4"\nDEL k3\nHMSET k3 "f" "5"\nDEL k1\nDEL k2\nHMSET k1 "a" 1 "b" 2 "c" 3\nHMSET k2 "d" "4" "e" "4"\nDEL k3\nHMSET k3 "f" "5"`)
		_ = QuantizeRedis(span)
	}
}
