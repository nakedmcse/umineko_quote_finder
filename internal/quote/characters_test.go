package quote

import (
	"testing"
)

func TestGetCharacterName_AllEntries(t *testing.T) {
	expect := map[string]string{
		"00":       "Group Voices",
		"01":       "Ushiromiya Kinzo",
		"02":       "Ushiromiya Krauss",
		"03":       "Ushiromiya Natsuhi",
		"04":       "Ushiromiya Jessica",
		"05":       "Ushiromiya Eva",
		"06":       "Ushiromiya Hideyoshi",
		"07":       "Ushiromiya George",
		"08":       "Ushiromiya Rudolf",
		"09":       "Ushiromiya Kyrie",
		"10":       "Ushiromiya Battler",
		"11":       "Ushiromiya Ange",
		"12":       "Ushiromiya Rosa",
		"13":       "Ushiromiya Maria",
		"14":       "Ronoue Genji",
		"15":       "Shannon",
		"16":       "Kanon",
		"17":       "Gohda Toshiro",
		"18":       "Kumasawa Chiyo",
		"19":       "Nanjo Terumasa",
		"20":       "Amakusa Juuza",
		"21":       "Okonogi Tetsuro",
		"22":       "Sumadera Kasumi",
		"23":       "Professor Ootsuki",
		"24":       "Captain Kawabata",
		"25":       "Nanjo Masayuki",
		"26":       "Kumasawa Sabakichi",
		"27":       "Beatrice",
		"28":       "Bernkastel",
		"29":       "Lambdadelta",
		"30":       "Virgilia",
		"31":       "Ronove",
		"32":       "Gaap",
		"33":       "Sakutarou",
		"34":       "Eva Beatrice",
		"35":       "Chiester 45",
		"36":       "Chiester 410",
		"37":       "Chiester 00",
		"38":       "Lucifer",
		"39":       "Leviathan",
		"40":       "Satan",
		"41":       "Belphegor",
		"42":       "Mammon",
		"43":       "Beelzebub",
		"44":       "Asmodeus",
		"45":       "Goat",
		"46":       "Furudo Erika",
		"47":       "Dlanor A. Knox",
		"48":       "Gertrude",
		"49":       "Cornelia",
		"50":       "Featherine",
		"51":       "Zepar",
		"52":       "Furfur",
		"53":       "Ushiromiya Lion",
		"54":       "Willard H. Wright",
		"55":       "Clair",
		"56":       "Hachijo Ikuko",
		"57":       "Hachijo Tohya",
		"58":       "Ushiromiya Kinzo",
		"59":       "Bice",
		"60":       "Beato the Elder",
		"99":       "Misc Voices",
		"narrator": "Narrator",
	}

	for id, wantName := range expect {
		got := CharacterNames.GetCharacterName(id)
		if got != wantName {
			t.Errorf("GetCharacterName(%q): got %q, want %q", id, got, wantName)
		}
	}

	if len(CharacterNames) != len(expect) {
		t.Errorf("CharacterNames has %d entries, test expects %d — update the test", len(CharacterNames), len(expect))
	}
}

func TestGetCharacterName_UnknownID(t *testing.T) {
	unknowns := []string{"", "100", "abc", "-1", "61"}
	for _, id := range unknowns {
		got := CharacterNames.GetCharacterName(id)
		if got != "Unknown" {
			t.Errorf("GetCharacterName(%q): got %q, want \"Unknown\"", id, got)
		}
	}
}

func TestGetAllCharacters_ReturnsCopy(t *testing.T) {
	all := CharacterNames.GetAllCharacters()

	if len(all) != len(CharacterNames) {
		t.Fatalf("GetAllCharacters returned %d entries, want %d", len(all), len(CharacterNames))
	}
	for id, wantName := range CharacterNames {
		if all[id] != wantName {
			t.Errorf("GetAllCharacters()[%q]: got %q, want %q", id, all[id], wantName)
		}
	}

	all["test_mutation"] = "should not leak"
	if _, exists := CharacterNames["test_mutation"]; exists {
		t.Error("GetAllCharacters did not return a copy — mutation leaked into CharacterNames")
	}
}
