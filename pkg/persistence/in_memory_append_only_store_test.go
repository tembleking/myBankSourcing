package persistence_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
)

var _ = Describe("InMemory / InMemoryAppendOnlyStore", func() {
	var (
		ctx   context.Context
		store *persistence.InMemoryAppendOnlyStore
	)

	BeforeEach(func() {
		ctx = context.Background()
		store = persistence.NewInMemoryAppendOnlyStore()
	})

	It("should be able to append to an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data")})

		Expect(err).To(BeNil())
	})

	It("should be able to append to multiple event streams", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data")})
		Expect(err).To(BeNil())

		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data")})
		Expect(err).To(BeNil())
	})

	It("should be able to append to multiple events to an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data1")})
		Expect(err).To(BeNil())

		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 1}, EventName: "eventName", EventData: []byte("data2")})
		Expect(err).To(BeNil())
	})

	When("there is a double append with the same expected version", func() {
		It("should return an error", func() {
			err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data")})
			Expect(err).To(BeNil())

			err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data")})
			Expect(err).To(MatchError(&persistence.ErrUnexpectedVersion{StreamName: "aggregate-0", Expected: 0}))
		})
	})

	It("should be able to read from an event stream", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data")})
		Expect(err).To(BeNil())

		data, err := store.ReadRecords(ctx, "aggregate-0")
		Expect(err).To(BeNil())
		Expect(data).To(HaveLen(1))
		Expect(data[0].ID.StreamName).To(Equal(persistence.StreamName("aggregate-0")))
		Expect(data[0].ID.StreamVersion).To(Equal(persistence.StreamVersion(0)))
		Expect(data[0].EventName).To(Equal("eventName"))
		Expect(data[0].EventData).To(Equal([]byte("data")))
	})

	It("should be able to retrieve the events by name", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data0")})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventName: "eventNameToIgnore", EventData: []byte("data1")})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data2")})
		Expect(err).To(BeNil())

		data, err := store.ReadEventsByName(ctx, "eventName")
		Expect(err).To(BeNil())
		Expect(data).To(HaveLen(2))
		Expect(data).To(ConsistOf(
			persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data0")},
			persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data2")},
		))
	})

	It("should be able to load all events", func() {
		err := store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data0")})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventName: "eventNameToIgnore", EventData: []byte("data1")})
		Expect(err).To(BeNil())
		err = store.Append(ctx, persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data2")})
		Expect(err).To(BeNil())

		data, err := store.ReadAllRecords(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(HaveLen(3))
		Expect(data).To(ConsistOf(
			persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data0")},
			persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventName: "eventNameToIgnore", EventData: []byte("data1")},
			persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data2")},
		))
	})

	When("there are no events", func() {
		It("doesn't return any event", func() {
			data, err := store.ReadRecords(ctx, "aggregate-0")
			Expect(err).ToNot(HaveOccurred())
			Expect(data).To(BeEmpty())

			data, err = store.ReadEventsByName(ctx, "eventName")
			Expect(err).ToNot(HaveOccurred())
			Expect(data).To(BeEmpty())
		})
	})

	Context("event dispatching", func() {
		When("there were no events pushed", func() {
			It("returns no events to dispatch", func() {
				undispatchedRecords, err := store.ReadUndispatchedRecords(ctx)

				Expect(err).ToNot(HaveOccurred())
				Expect(undispatchedRecords).To(BeEmpty())
			})
		})

		When("there are events pushed and were not marked as dispatched", func() {
			It("returns the undispatched events", func() {
				event := persistence.StoredStreamEvent{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data0")}
				err := store.Append(ctx, event)
				Expect(err).ToNot(HaveOccurred())

				records, err := store.ReadUndispatchedRecords(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(records).To(ConsistOf(event))
			})
		})

		When("the event was marked as dispatched", func() {
			It("is not returned again", func() {
				err := store.Append(ctx, eventsStored()...)
				Expect(err).ToNot(HaveOccurred())

				err = store.MarkRecordsAsDispatched(ctx, eventsStoredIDs()...)
				Expect(err).ToNot(HaveOccurred())

				records, err := store.ReadUndispatchedRecords(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(records).To(BeEmpty())
			})
		})

		When("the undispatched events are returned", func() {
			It("doesn't return them again before some time passes", func() {
				err := store.Append(ctx, eventsStored()...)
				Expect(err).ToNot(HaveOccurred())

				records, err := store.ReadUndispatchedRecords(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(records).To(ConsistOf(eventsStored()))

				records, err = store.ReadUndispatchedRecords(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(records).To(BeEmpty())
			})

			When("some time passes", func() {
				It("returns the undispatched events again", func() {
					err := store.Append(ctx, eventsStored()...)
					Expect(err).ToNot(HaveOccurred())

					records, err := store.ReadUndispatchedRecords(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(records).To(ConsistOf(eventsStored()))

					time.Sleep(6 * time.Second)

					records, err = store.ReadUndispatchedRecords(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(records).To(ConsistOf(eventsStored()))
				})

				When("the event was marked as dispatched", func() {
					It("is not returned again", func() {
						err := store.Append(ctx, eventsStored()...)
						Expect(err).ToNot(HaveOccurred())

						records, err := store.ReadUndispatchedRecords(ctx)
						Expect(err).ToNot(HaveOccurred())
						Expect(records).To(ConsistOf(eventsStored()))

						err = store.MarkRecordsAsDispatched(ctx, eventsStoredIDs()...)
						Expect(err).ToNot(HaveOccurred())

						time.Sleep(6 * time.Second)

						records, err = store.ReadUndispatchedRecords(ctx)
						Expect(err).ToNot(HaveOccurred())
						Expect(records).To(BeEmpty())
					})
				})
			})
		})
	})
})

func eventsStored() []persistence.StoredStreamEvent {
	return []persistence.StoredStreamEvent{
		{ID: persistence.StreamID{StreamName: "aggregate-0", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data0")},
		{ID: persistence.StreamID{StreamName: "aggregate-1", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data1")},
		{ID: persistence.StreamID{StreamName: "aggregate-2", StreamVersion: 0}, EventName: "eventName", EventData: []byte("data2")},
	}
}

func eventsStoredIDs() []persistence.StreamID {
	var ids []persistence.StreamID
	for _, event := range eventsStored() {
		ids = append(ids, event.ID)
	}
	return ids
}
