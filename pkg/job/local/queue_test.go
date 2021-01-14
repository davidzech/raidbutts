package local

import (
	"testing"

	"github.com/davidzech/raidbutts/pkg/job"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := NewQueue()

	var job job.Job
	assert.NoError(t, q.Publish(&job))

	out, err := q.Consume()
	assert.NoError(t, err)
	assert.Equal(t, &job, out)
}
