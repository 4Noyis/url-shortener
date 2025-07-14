package filter

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type BloomFilter struct {
	filter *bloom.BloomFilter
}

// NewBloomFilter creates a new Bloom filter with specified capacity and false positive rate
func NewBloomFilter(expectedElements uint, falsePositiveRate float64) *BloomFilter {
	filter := bloom.NewWithEstimates(expectedElements, falsePositiveRate)
	return &BloomFilter{
		filter: filter,
	}
}

// Add inserts an item into the Bloom filter
func (bf *BloomFilter) Add(item string) {
	bf.filter.Add([]byte(item))
}

// Test checks if an item might exist in the Bloom filter (may have false positives)
func (bf *BloomFilter) Test(item string) bool {
	return bf.filter.Test([]byte(item))
}

// AddAll adds multiple items to the Bloom filter in batch
func (bf *BloomFilter) AddAll(items []string) {
	for _, item := range items {
		bf.Add(item)
	}
}