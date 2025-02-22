package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"strings"
)

func Compress_str(data string) string {
	var buff bytes.Buffer
	writer := gzip.NewWriter(&buff)

	_, err := writer.Write([]byte(data))
	if err != nil {
		fmt.Println("ERROR IN COMPRESSING DATA : ", err)
	}
	writer.Close()

	return buff.String()
}

func Decompress_str(data string) (string, error) {
	reader, err := gzip.NewReader(bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("ERROR IN DECOMPRESSING: ", err)
		return "", err
	}
	decompressed_bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("error in decompression: ", err)
		return "", err
	}

	return string(decompressed_bytes), nil
}

// Fonction pour insérer une valeur entre /start et /end
func InsertQuery(chaine, valeur string) string {
	start := "/start"
	end := "/end"

	// Vérifie si les balises existent déjà
	if !strings.Contains(chaine, start) || !strings.Contains(chaine, end) {
		return chaine + " " + start + valeur + end // Ajoute les balises avec la valeur
	}

	// Retire l'ancienne valeur si elle existe
	chaine = retirerValeur(chaine)

	// Ajoute la nouvelle valeur
	return fmt.Sprintf("%s %s%s%s", chaine[:strings.Index(chaine, end)], start, valeur, end)
}

// Fonction pour extraire la valeur entre /start et /end
func ExtractQuery(chaine string) string {
	start := "/start"
	end := "/end"

	startIndex := strings.Index(chaine, start)
	endIndex := strings.Index(chaine, end)

	if startIndex == -1 || endIndex == -1 || endIndex < startIndex {
		return "" // Retourne une chaîne vide si les balises n'existent pas
	}

	return chaine[startIndex+len(start) : endIndex]
}

// Fonction pour retirer la valeur entre /start et /end
func retirerValeur(chaine string) string {
	start := "/start"
	end := "/end"

	if strings.Contains(chaine, start) && strings.Contains(chaine, end) {
		return strings.Replace(chaine, chaineBetween(chaine, start, end), "", 1)
	}
	return chaine // Retourne la chaîne sans modification si les balises n'existent pas
}

func GetId(chaine string) (string, string) {
	n := 8
	if len(chaine) < n {
		return "", ""
	}
	return chaine[:n], chaine[n+1:]
}

// Fonction pour extraire le texte entre /start et /end
func chaineBetween(chaine, start, end string) string {
	startIndex := strings.Index(chaine, start)
	endIndex := strings.Index(chaine, end)

	if startIndex == -1 || endIndex == -1 || endIndex < startIndex {
		return ""
	}

	return chaine[startIndex+len(start) : endIndex]
}

/*func main() {
	// Exemple d'utilisation
	chaine := "Ceci est un exemple."

	fmt.Println("Chaîne initiale : ", chaine)

	// Insérer une nouvelle valeur
	chaine = insererValeur(chaine, "NouvelleValeur")
	fmt.Println("Après insertion : ", chaine)

	// Extraire la valeur
	valeurExtraite := extraireValeur(chaine)
	fmt.Println("Valeur extraite : ", valeurExtraite)

	// Retirer la valeur
	chaine = retirerValeur(chaine)
	fmt.Println("Après retrait : ", chaine)
}
*/
