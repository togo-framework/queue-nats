package queuenats

import "testing"

func TestSubject(t *testing.T) {
	q := &natsQueue{prefix: "togo.queue"}
	if got := q.subject("emails"); got != "togo.queue.emails" {
		t.Fatalf("subject = %q", got)
	}
}
