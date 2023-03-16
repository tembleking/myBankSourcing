package inmemory_test

import (
	"context"
	"sort"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tembleking/myBankSourcing/pkg/persistence"
	"github.com/tembleking/myBankSourcing/pkg/persistence/inmemory"
)

var _ = Describe("InMemory / AppendOnlyStore", func() {
	var (
		store *inmemory.AppendOnlyStore
	)

	BeforeEach(func() {
		store = inmemory.NewAppendOnlyStore()
	})

	It("should be able to append to an event stream", func() {
		err := store.Append(context.Background(), "aggregate-0", []byte("data"), 0)

		Expect(err).To(BeNil())
	})

	It("should be able to append to multiple event streams", func() {
		err := store.Append(context.Background(), "aggregate-0", []byte("data"), 0)
		Expect(err).To(BeNil())

		err = store.Append(context.Background(), "aggregate-1", []byte("data"), 0)
		Expect(err).To(BeNil())
	})

	It("should be able to append to multiple events to an event stream", func() {
		err := store.Append(context.Background(), "aggregate-0", []byte("data"), 0)
		Expect(err).To(BeNil())

		err = store.Append(context.Background(), "aggregate-0", []byte("data"), 1)
		Expect(err).To(BeNil())
	})

	When("there is a double append with the same expected version", func() {
		It("should return an error", func() {
			err := store.Append(context.Background(), "aggregate-0", []byte("data"), 0)
			Expect(err).To(BeNil())

			err = store.Append(context.Background(), "aggregate-0", []byte("data"), 0)
			Expect(err).To(MatchError(&persistence.ErrUnexpectedVersion{Found: 1, Expected: 0}))
		})
	})

	When("the expected version is not met", func() {
		It("should return an error", func() {
			err := store.Append(context.Background(), "aggregate-0", []byte("data"), 1)

			Expect(err).To(MatchError(&persistence.ErrUnexpectedVersion{Found: 0, Expected: 1}))
		})
	})

	It("should be able to read from an event stream", func() {
		err := store.Append(context.Background(), "aggregate-0", []byte("data"), 0)
		Expect(err).To(BeNil())

		data, err := store.ReadRecords(context.Background(), "aggregate-0", 0, 0)
		Expect(err).To(BeNil())
		Expect(data).To(HaveLen(1))
		Expect(data[0].Data).To(Equal([]byte("data")))
		Expect(data[0].Version).To(Equal(uint64(1)))
	})

	When("the requested start is 1", func() {
		It("should return all events but the previous ones", func() {
			err := store.Append(context.Background(), "aggregate-0", []byte("data-0"), 0)
			Expect(err).To(BeNil())

			err = store.Append(context.Background(), "aggregate-0", []byte("data-1"), 1)
			Expect(err).To(BeNil())

			data, err := store.ReadRecords(context.Background(), "aggregate-0", 1, 0)
			Expect(err).To(BeNil())
			Expect(data).To(HaveLen(1))
			Expect(data[0].Data).To(Equal([]byte("data-1")))
			Expect(data[0].Version).To(Equal(uint64(2)))
		})
	})

	When("the start version is higher than the number of events", func() {
		It("should return an empty list", func() {
			err := store.Append(context.Background(), "aggregate-0", []byte("data"), 0)
			Expect(err).To(BeNil())

			data, err := store.ReadRecords(context.Background(), "aggregate-0", 1, 0)
			Expect(err).To(BeNil())
			Expect(data).To(HaveLen(0))
		})

		When("the max count is set", func() {
			It("should return an empty list", func() {
				err := store.Append(context.Background(), "aggregate-0", []byte("data"), 0)
				Expect(err).To(BeNil())

				data, err := store.ReadRecords(context.Background(), "aggregate-0", 1, 1)
				Expect(err).To(BeNil())
				Expect(data).To(HaveLen(0))
			})
		})
	})

	When("the max count is set", func() {
		It("should return the max count of events", func() {
			err := store.Append(context.Background(), "aggregate-0", []byte("data-0"), 0)
			Expect(err).To(BeNil())

			err = store.Append(context.Background(), "aggregate-0", []byte("data-1"), 1)
			Expect(err).To(BeNil())

			data, err := store.ReadRecords(context.Background(), "aggregate-0", 0, 1)
			Expect(err).To(BeNil())
			Expect(data).To(HaveLen(1))
			Expect(data[0].Data).To(Equal([]byte("data-0")))
		})
	})

	It("should be able to read from all event streams", func() {
		err := store.Append(context.Background(), "aggregate-0", []byte("data0"), 0)
		Expect(err).To(BeNil())

		err = store.Append(context.Background(), "aggregate-1", []byte("data1"), 0)
		Expect(err).To(BeNil())

		data, err := store.ReadAllRecords(context.Background(), 0, 0)
		Expect(err).To(BeNil())
		Expect(data).To(HaveLen(2))
		// sort for testing only, we don't care about the order in production
		sort.Slice(data, func(i, j int) bool { return data[i].Name < data[j].Name })
		Expect(data[0]).To(Equal(persistence.DataWithName{
			Name: "aggregate-0",
			Data: []byte("data0"),
		}))
		Expect(data[1]).To(Equal(persistence.DataWithName{
			Name: "aggregate-1",
			Data: []byte("data1"),
		}))
	})
})
