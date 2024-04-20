package benching

import (
	"myapp/pkg/database"
	"myapp/pkg/words"
	"testing"
)

func normal_query(query string) []string {
	arr_query := words.SplitString(query)
	normal_query, _ := words.Stemming(arr_query)
	return normal_query
}

func BenchmarkSearch(b *testing.B) {
	arr := normal_query("I'm following your questions")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		database.SearchDatabase(arr)
	}

}

func BenchmarkIndexedSearch(b *testing.B) {
	arr := normal_query("I'm following your questions")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		database.SearchByIndex(arr)
	}
}
