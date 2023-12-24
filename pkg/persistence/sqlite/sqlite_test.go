package sqlite_test

import (
	"context"
	"log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/sqlite"
)

var _ = Describe("Sqlite AppendOnlyStore", func() {
	var (
		ctx   context.Context
		store *sqlite.AppendOnlyStore
	)

	BeforeEach(func() {
		log.SetOutput(GinkgoWriter)
		ctx = context.Background()
		store = setupStore()
	})

	AfterEach(func() {
		store.Close()
	})

	It("should be able to append to an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventID: "event0", EventName: "eventName", EventData: []byte("data"), ContentType: "some-content-type"})
		Expect(err).To(BeNil())
	})

	It("should be able to append to multiple event streams", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventID: "event0", EventName: "eventName", EventData: []byte("data"), ContentType: "some-content-type"})
		Expect(err).To(BeNil())

		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventID: "event1", EventName: "eventName", EventData: []byte("data"), ContentType: "some-content-type"})
		Expect(err).To(BeNil())
	})

	It("should be able to append to multiple events to an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventID: "event0", EventName: "eventName", EventData: []byte("data1"), ContentType: "some-content-type"})
		Expect(err).To(BeNil())

		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventID: "event1", EventName: "eventName", EventData: []byte("data2"), ContentType: "some-content-type"})
		Expect(err).To(BeNil())
	})

	When("there is a double append with the same expected version", func() {
		It("should return an error", func() {
			err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data"), ContentType: "some-content-type"})
			Expect(err).To(BeNil())

			err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data"), ContentType: "some-content-type"})
			Expect(err).To(MatchError(persistence.ErrUnexpectedVersion))
		})
	})

	It("should be able to read from an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data"), ContentType: "some-content-type"})
		Expect(err).To(BeNil())

		data, err := store.ReadRecords(ctx, "aggregate-0")
		Expect(err).To(BeNil())
		Expect(data).To(HaveLen(1))
		Expect(data[0].ID.StreamName).To(Equal("aggregate-0"))
		Expect(data[0].EventName).To(Equal("eventName"))
		Expect(data[0].EventData).To(Equal([]byte("data")))
		Expect(data[0].ID.StreamVersion).To(Equal(uint64(0)))
	})

	It("should return all the events", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventID: "event0", EventName: "eventName", EventData: []byte("data0"), ContentType: "some-content-type-0"})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventID: "event1", EventName: "eventNameToIgnore", EventData: []byte("data1"), ContentType: "some-content-type-1"})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventID: "event2", EventName: "eventName", EventData: []byte("data2-0"), ContentType: "some-content-type-0"})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 1}, EventID: "event3", EventName: "eventName", EventData: []byte("data2-1"), ContentType: "some-content-type-1"})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 2}, EventID: "event4", EventName: "eventName", EventData: []byte("data2-2"), ContentType: "some-content-type-2"})
		Expect(err).To(BeNil())

		data, err := store.ReadAllRecords(ctx)
		Expect(err).To(BeNil())
		Expect(data).To(HaveLen(5))
		Expect(data).To(Equal([]persistence.StoredStreamEvent{
			{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventID: "event0", EventName: "eventName", EventData: []byte("data0"), ContentType: "some-content-type-0"},
			{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventID: "event1", EventName: "eventNameToIgnore", EventData: []byte("data1"), ContentType: "some-content-type-1"},
			{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventID: "event2", EventName: "eventName", EventData: []byte("data2-0"), ContentType: "some-content-type-0"},
			{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 1}, EventID: "event3", EventName: "eventName", EventData: []byte("data2-1"), ContentType: "some-content-type-1"},
			{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 2}, EventID: "event4", EventName: "eventName", EventData: []byte("data2-2"), ContentType: "some-content-type-2"},
		}))
	})

	When("filtering events after other eventID", func() {
		BeforeEach(func() {
			Expect(store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventID: "event0", EventName: "eventName", EventData: []byte("data0"), ContentType: "some-content-type-0"})).To(Succeed())
			Expect(store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventID: "event1", EventName: "eventNameToIgnore", EventData: []byte("data1"), ContentType: "some-content-type-1"})).To(Succeed())
			Expect(store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventID: "event2", EventName: "eventName", EventData: []byte("data2-0"), ContentType: "some-content-type-0"})).To(Succeed())
			Expect(store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 1}, EventID: "event3", EventName: "eventName", EventData: []byte("data2-1"), ContentType: "some-content-type-1"})).To(Succeed())
			Expect(store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 2}, EventID: "event4", EventName: "eventName", EventData: []byte("data2-2"), ContentType: "some-content-type-2"})).To(Succeed())
			Expect(store.ReadAllRecords(ctx)).To(HaveLen(5))
		})

		AfterEach(func() {
			// Checks if the store is not modified by the filter methods.
			Expect(store.ReadAllRecords(ctx)).To(HaveLen(5))
		})

		It("should return only the events after the eventID", func() {
			records, err := store.AfterEventID("event2").ReadAllRecords(ctx)

			Expect(err).To(BeNil())
			Expect(records).To(Equal([]persistence.StoredStreamEvent{
				{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 1}, EventID: "event3", EventName: "eventName", EventData: []byte("data2-1"), ContentType: "some-content-type-1"},
				{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 2}, EventID: "event4", EventName: "eventName", EventData: []byte("data2-2"), ContentType: "some-content-type-2"},
			}))
		})

		It("should return a limit number of events", func() {
			records, err := store.Limit(2).ReadAllRecords(ctx)

			Expect(err).To(BeNil())
			Expect(records).To(Equal([]persistence.StoredStreamEvent{
				{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventID: "event0", EventName: "eventName", EventData: []byte("data0"), ContentType: "some-content-type-0"},
				{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventID: "event1", EventName: "eventNameToIgnore", EventData: []byte("data1"), ContentType: "some-content-type-1"},
			}))
		})

		It("should return a limit number of events after the eventID", func() {
			records, err := store.AfterEventID("event1").Limit(2).ReadAllRecords(ctx)

			Expect(err).To(BeNil())
			Expect(records).To(Equal([]persistence.StoredStreamEvent{
				{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventID: "event2", EventName: "eventName", EventData: []byte("data2-0"), ContentType: "some-content-type-0"},
				{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 1}, EventID: "event3", EventName: "eventName", EventData: []byte("data2-1"), ContentType: "some-content-type-1"},
			}))
		})
	})
})

func setupStore() *sqlite.AppendOnlyStore {
	db := sqlite.InMemory()

	err := db.MigrateDB(context.Background())
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return db
}
