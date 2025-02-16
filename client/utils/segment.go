package utils

func segmentQuery(chaine string, tailleSegment int) []string {
	var segments []string
	// Vérifie si la taille de la chaîne est supérieure à 5
	if len(chaine) <= tailleSegment {
		return []string{chaine} // Retourne la chaîne entière si elle est trop courte
	}

	for i := 0; i < len(chaine); i += tailleSegment {
		// Calcule la fin du segment
		fin := i + tailleSegment
		if fin > len(chaine) {
			fin = len(chaine) // Évite de dépasser la longueur de la chaîne
		}
		segments = append(segments, chaine[i:fin])
	}

	return segments
}
