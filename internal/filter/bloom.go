package filter

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type BloomFilter struct {
	filter *bloom.BloomFilter
}

func NewBloomFilter(expectedElements uint, falsePositiveRate float64) *BloomFilter {
	filter := bloom.NewWithEstimates(expectedElements, falsePositiveRate)
	return &BloomFilter{
		filter: filter,
	}
}

func (bf *BloomFilter) Add(item string) {
	bf.filter.Add([]byte(item))
}

func (bf *BloomFilter) Test(item string) bool {
	return bf.filter.Test([]byte(item))
}

func (bf *BloomFilter) AddAll(items []string) {
	for _, item := range items {
		bf.Add(item)
	}
}