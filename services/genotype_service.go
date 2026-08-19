package services

import "strings"

type GenotypeService struct{}

func NewGenotypeService() *GenotypeService {
	return &GenotypeService{}
}

// IsCompatible checks if two genotypes can safely marry without risk of Sickle Cell Disease (SS/SC) in offspring.
// Medical rules:
// - AA is compatible with AA, AS, AC, SS, SC
// - AS is compatible ONLY with AA
// - AC is compatible ONLY with AA
// - SS / SC are compatible ONLY with AA
func (gs *GenotypeService) IsCompatible(g1, g2 string) bool {
	g1 = strings.ToUpper(strings.TrimSpace(g1))
	g2 = strings.ToUpper(strings.TrimSpace(g2))

	// If either genotype is missing/unknown, we do not strictly block, but return unverified
	if g1 == "" || g2 == "" || g1 == "UNKNOWN" || g2 == "UNKNOWN" {
		return true
	}

	if g1 == "AA" || g2 == "AA" {
		return true
	}

	// Any non-AA + non-AA pairing carries high medical risk (AS+AS = 25% SS; AS+AC = 25% SC; AS+SS = 50% SS)
	return false
}

// EvaluateCompatibility returns dynamic status ("compatible", "incompatible", "unverified") and medical advisory text.
func (gs *GenotypeService) EvaluateCompatibility(userGenotype, candidateGenotype string) (string, string) {
	g1 := strings.ToUpper(strings.TrimSpace(userGenotype))
	g2 := strings.ToUpper(strings.TrimSpace(candidateGenotype))

	if g1 == "" || g2 == "" || g1 == "UNKNOWN" || g2 == "UNKNOWN" {
		return "unverified", "Genotype unverified. Verification recommended before serious commitment."
	}

	if gs.IsCompatible(g1, g2) {
		return "compatible", "Genotype Compatible: No Sickle Cell risk for offspring."
	}

	// Detail specific medical warnings
	if (g1 == "AS" && g2 == "AS") || (g1 == "AS" && g2 == "SS") || (g1 == "SS" && g2 == "AS") || (g1 == "SS" && g2 == "SS") {
		return "incompatible", "High Medical Risk: AS/SS pairing carries 25% to 50% risk of Sickle Cell Anemia (SS) in children."
	}
	if (g1 == "AS" && g2 == "AC") || (g1 == "AC" && g2 == "AS") || (g1 == "AC" && g2 == "AC") || (g1 == "AC" && g2 == "SC") {
		return "incompatible", "High Medical Risk: Carrier pairing carries 25% risk of Sickle Cell Hemoglobinopathy (SC/CC) in children."
	}

	return "incompatible", "Medical Risk: Incompatible genotypes carry significant risk of hemoglobinopathy in offspring."
}
