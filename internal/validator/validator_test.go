package validator_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/nestoroprysk/mood-tracker/internal/registry/registryv1"
	"github.com/nestoroprysk/mood-tracker/internal/validator"
)

func TestValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validator Suite")
}

var _ = DescribeTable("Validates registryv1.Item", func(t registryv1.Item, succeed bool) {
	v := validator.New()
	result := v.Struct(t)
	if succeed {
		Expect(result).To(Succeed())
	} else {
		Expect(result).NotTo(BeNil())
	}
},
	Entry("Errors if invalid item",
		registryv1.Item{},
		false, /* succeed */
	),
	Entry("Errors if item without mood",
		registryv1.Item{Time: time.Now()},
		false, /* succeed */
	),
	Entry("Errors if too good mood",
		registryv1.Item{
			Time: time.Now(),
			Mood: 6,
		},
		false, /* succeed */
	),
	Entry("Errors if too bad mood",
		registryv1.Item{
			Time: time.Now(),
			Mood: -1,
		},
		false, /* succeed */
	),
	Entry("Accepts good mood",
		registryv1.Item{
			Time: time.Now(),
			Mood: 5,
		},
		true, /* succeed */
	),
	Entry("Accepts bad mood",
		registryv1.Item{
			Time: time.Now(),
			Mood: 1,
		},
		true, /* succeed */
	),
	Entry("Accepts a label",
		registryv1.Item{
			Time:   time.Now(),
			Mood:   1,
			Labels: []string{"happy"},
		},
		true, /* succeed */
	),
	Entry("Rejects a short label",
		registryv1.Item{
			Time:   time.Now(),
			Mood:   1,
			Labels: []string{"up"},
		},
		false, /* succeed */
	),
)
