package surrealdb_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	surreal "github.com/surrealdb/surrealdb.go"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/surrealdb"
)

var _ = Describe("SurrealDB AppendOnlyStore", func() {
	var (
		ctx   context.Context
		store *surrealdb.AppendOnlyStore
	)

	BeforeEach(func() {
		ctx = context.Background()
		db, err := surreal.New("ws://localhost:8000/rpc")
		Expect(err).ToNot(HaveOccurred())
		_, err = db.Signin(struct {
			User string `json:"user"`
			Pass string `json:"pass"`
		}{
			User: "root",
			Pass: "root",
		})
		Expect(err).ToNot(HaveOccurred())
		_, err = db.Use("ns", "db")
		Expect(err).ToNot(HaveOccurred())
		_, err = db.Delete("event")
		Expect(err).ToNot(HaveOccurred())

		store = surrealdb.NewAppendOnlyStore(db)

	})

	It("should be able to append to an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 0, EventName: "eventName", EventData: []byte("data")})
		Expect(err).To(BeNil())
	})

	It("should be able to append to multiple event streams", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 0, EventName: "eventName", EventData: []byte("data")})
		Expect(err).To(BeNil())

		err = store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-1", StreamVersion: 0, EventName: "eventName", EventData: []byte("data")})
		Expect(err).To(BeNil())
	})

	It("should be able to append to multiple events to an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 0, EventName: "eventName", EventData: []byte("data1")})
		Expect(err).To(BeNil())

		err = store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 1, EventName: "eventName", EventData: []byte("data2")})
		Expect(err).To(BeNil())
	})

	When("there is a double append with the same expected version", func() {
		It("should return an error", func() {
			err := store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 0, EventName: "eventName", EventData: []byte("data")})
			Expect(err).To(BeNil())

			err = store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 0, EventName: "eventName", EventData: []byte("data")})
			Expect(err).To(MatchError(&persistence.ErrUnexpectedVersion{Found: 1, Expected: 0}))
		})
	})

	When("the expected version is not met", func() {
		It("should return an error", func() {
			err := store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 1, EventName: "eventName", EventData: []byte("data")})

			Expect(err).To(MatchError(&persistence.ErrUnexpectedVersion{Found: 0, Expected: 1}))
		})
	})

	It("should be able to read from an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 0, EventName: "eventName", EventData: []byte("data")})
		Expect(err).To(BeNil())

		data, err := store.ReadRecords(ctx, "aggregate-0")
		Expect(err).To(BeNil())
		Expect(data).To(HaveLen(1))
		Expect(data[0].StreamID).To(Equal("aggregate-0"))
		Expect(data[0].EventName).To(Equal("eventName"))
		Expect(data[0].EventData).To(Equal([]byte("data")))
		Expect(data[0].StreamVersion).To(Equal(uint64(0)))
	})

	It("should be able to retrieve the events by name", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 0, EventName: "eventName", EventData: []byte("data0")})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-1", StreamVersion: 0, EventName: "eventNameToIgnore", EventData: []byte("data1")})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-2", StreamVersion: 0, EventName: "eventName", EventData: []byte("data2-0")})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-2", StreamVersion: 1, EventName: "eventName", EventData: []byte("data2-1")})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{StreamID: "aggregate-2", StreamVersion: 2, EventName: "eventName", EventData: []byte("data2-2")})
		Expect(err).To(BeNil())

		data, err := store.ReadEventsByName(ctx, "eventName")
		Expect(err).To(BeNil())
		Expect(data).To(HaveLen(4))
		Expect(data).To(ConsistOf(
			persistence.StoredStreamEvent{StreamID: "aggregate-0", StreamVersion: 0, EventName: "eventName", EventData: []byte("data0")},
			persistence.StoredStreamEvent{StreamID: "aggregate-2", StreamVersion: 0, EventName: "eventName", EventData: []byte("data2-0")},
			persistence.StoredStreamEvent{StreamID: "aggregate-2", StreamVersion: 1, EventName: "eventName", EventData: []byte("data2-1")},
			persistence.StoredStreamEvent{StreamID: "aggregate-2", StreamVersion: 2, EventName: "eventName", EventData: []byte("data2-2")},
		))
	})
})
