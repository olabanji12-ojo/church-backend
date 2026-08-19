package main

import (
	"fmt"
	"github.com/olabanji12-ojo/church-backend/services"
)

func main() {
	gs := services.NewGenotypeService()

	tests := []struct {
		g1, g2   string
		expected bool
	}{
		{"AA", "AA", true},
		{"AA", "AS", true},
		{"AA", "AC", true},
		{"AA", "SS", true},
		{"AS", "AA", true},
		{"AS", "AS", false},
		{"AS", "AC", false},
		{"AS", "SS", false},
		{"AC", "AC", false},
		{"SS", "SS", false},
		{"AS", "Unknown", true},
	}

	fmt.Println("=== RUNNING GENOTYPE COMPATIBILITY SERVICE TESTS ===")
	passed := 0
	for i, t := range tests {
		result := gs.IsCompatible(t.g1, t.g2)
		status, warning := gs.EvaluateCompatibility(t.g1, t.g2)
		if result == t.expected {
			fmt.Printf("✅ Test %d Passed: %s + %s -> Compatible: %v | Status: %s | Warning: %s\n", i+1, t.g1, t.g2, result, status, warning)
			passed++
		} else {
			fmt.Printf("❌ Test %d Failed: %s + %s -> Got %v, Expected %v\n", i+1, t.g1, t.g2, result, t.expected)
		}
	}

	fmt.Printf("\nPassed %d / %d tests successfully.\n", passed, len(tests))
}
