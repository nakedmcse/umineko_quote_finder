package quote

type characterMapping map[string]string

var CharacterNames = characterMapping{
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

func (c characterMapping) GetCharacterName(id string) string {
	if name, ok := c[id]; ok {
		return name
	}
	return "Unknown"
}

func (c characterMapping) GetAllCharacters() map[string]string {
	out := make(map[string]string, len(c))
	for k, v := range c {
		out[k] = v
	}
	return out
}
