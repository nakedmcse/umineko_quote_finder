package audio

import (
	"encoding/base64"
	"fmt"
	"os"
)

const halfASecondSilence = "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AJABAMAYg5BBBiFjEEIIIYQQQggJAAAYcAAACDChDBQashIASAUAAAxSijEHpaQUKcWYcxBKSSlSijHnIJSSUsWYcxBKSam1ijHnIJSSUmudcxBKSam1GDvnIJSSUmsxhhBKSam1GGMMIZSSUmsx1lpKSam1GGvMtZSSUmsx1lprSq21GGutNeeUWmsx1lpzzgUAIDQ4AIAd2LA6wknRWGChISsBgDwAAEgpxhhjjDGlFGOMMcaYUooxxhhjjDHGGGOMMaYYY4wxxhhjjDHGGGOMMcYYY4wxxhhjjDHGGGOMMcYYY8w5xhhjjDnmHGOMMcacc04AAFCBAwBAgI0imxOMBBUashIACAcAAIxhzDnnIJSQSqOUcxBCKCWVVhqlnIMSQikptZY5JyWlUlJqLbbMOSkplZJSay12ElJqLaXWYqyxg5BSa6m1FmONHYRSWootxhpz7SCUklprMcZaayilpdhirLHWmkMpqbUWY60151xSai3GWmvNteeSUmsxxlprrbmn1mKssdZcc+89tRZjjbXmnHvOBQCYPDgAQCXYOMNK0lnhaHChISsBgNwAAEYpxpxzDkIIIYQQQgiVUow55xyEEEIIIYQQKqUYc845CCGEEEIIIWSMOeccdBBCCCGEEELIGHPOOQghhBBCCCGE0DnnHIQQQgghhFBCKaVzzjkHIYQQQgghhFBK5xyEEEIIIYQSSiillM45CCGEEEIIpYRSSikhhBBCCCGEEkoppZRSOgghhBBCCKWUUkoppYQQQgghhBBKKaWUUkoJIYQQQgghlFJKKaWUEkIIIYQSSimllFJKKSWEEEIIoZRSSimllFJKCCGEUkoppZRSSimllBBCKCGUUkoppZRSSikhhBJKKKWUUkoppZRSQggllFJKKaWUUkoppYQQQiillFJKKaWUUkoJIZRSSimllFJKKaWUUgAA0IEDAECAEZUWYqcZVx6BIwoZJqBCQ1YCAOEAAAAhlFJKKaWUUmokpZRSSimllFIjJaWUUkoppZRSSimllFJKKaWUUkoppZRSSimllFJKKaWUUkoppZRSSimllFJKKaWUUkoppZRSSimllFJKKaWUUkoppZRSSimllFJKKaWUUkoppZRSSimllFJKKaWUUkoppZRSSqWUUkoppZRSSimllFJKKaWUUkoppQCoywwHwOgJG2dYSTorHA0uNGQlAJAWAAAYw5hjjkEnoZSUWmuYglBC6KSk0kpssTVKQQghhFJSSq211jLoqJSSSkqtxRZjjJmDUlIqJaXUYoyx1g5CSi21FluLseZaawehpJRaiy3GWmuuvYOQSmut5RhjsDnn2kEoKbXYYow111p7Dqm0FmOMtfZca805iFJSijHWGnPNNffcS0qtxZprrjUHn3MQpqXYao0154x7EDr41FqNueYedNBB5x50Sq3WWmvOPQchfPC5tVhrzTXn3oMPOgjfaqs151xr7z33noNuMdZcc9DBByF88EG4GGvPOfcchA46+B4MAMiNcABAXDCSkDrLsNKIG0/AEIEUGrIKAIgBACCMQQYhhJRSSimllGKKKcYYY4wxxhhjjDHGGGOMMcYEAAAmOAAABFjBrszSqo3ipk7yog8Cn9ARm5Ehl1IxkxNBj9RQi5Vgh1ZwgxeAhYasBADIAAAQiLHmWnOOEJTWYu25VEo5arHnlCGCnLScS8kMQU5aay1kyCgnMbYUMoQUtNpa6ZRSjGKrsXSMMUmpxZZK5yAAAACCAAADETITCBRAgYEMADhASJACAAoLDB3DRUBALiGjwKBwTDgnnTYAAEGIzBCJiMUgMaEaKCqmA4DFBYZ8AMjQ2Ei7uIAuA1zQxV0HQghCEIJYHEABCTg44YYn3vCEG5ygU1TqIAAAAAAAEADgAQAg2QAiopmZ4+jw+AAJERkhKTE5QUlRCQAAAAAAIAD4AABIVoCIaGbmODo8PkBCREZISkxOUFJUAgAAAAAAAAAAgICAAAAAAABAAAAAgIBPZ2dTAAQiVgAAAAAAAL4vAAACAAAAa5iZBhcBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQAKDg4ODg4ODg4ODg4ODg4ODg4ODg4O"

type FilePathResolver func(charId, audioId string) string

type Combiner interface {
	CombineOgg(charId string, ids []string, resolve FilePathResolver) ([]byte, error)
}

type combiner struct {
	silencePages []oggPage
}

func NewCombiner() (Combiner, error) {
	silenceOgg, err := base64.StdEncoding.DecodeString(halfASecondSilence)
	if err != nil {
		return nil, fmt.Errorf("failed to decode silence audio: %v", err)
	}
	pages, err := parseOggPages(silenceOgg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse silence audio: %v", err)
	}
	return &combiner{silencePages: pages}, nil
}

func (c *combiner) CombineOgg(charId string, ids []string, resolve FilePathResolver) ([]byte, error) {
	var allFilePages [][]oggPage

	for i := 0; i < len(ids); i++ {
		filePath := resolve(charId, ids[i])
		if filePath == "" {
			return nil, fmt.Errorf("audio file not found: %s", ids[i])
		}
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read audio file: %s", ids[i])
		}
		pages, err := parseOggPages(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse OGG file %s: %v", ids[i], err)
		}
		allFilePages = append(allFilePages, pages)
	}

	if len(allFilePages) == 0 {
		return nil, fmt.Errorf("no audio files to combine")
	}

	if len(allFilePages) > 1 {
		var withSilence [][]oggPage
		for i := 0; i < len(allFilePages); i++ {
			withSilence = append(withSilence, allFilePages[i])
			if i < len(allFilePages)-1 {
				withSilence = append(withSilence, c.silencePages)
			}
		}
		allFilePages = withSilence
	}

	serialNumber := allFilePages[0][0].serialNumber
	var result []byte
	var sequenceNum uint32
	var granuleOffset int64

	for fileIdx := 0; fileIdx < len(allFilePages); fileIdx++ {
		pages := allFilePages[fileIdx]
		isFirst := fileIdx == 0
		isLast := fileIdx == len(allFilePages)-1

		var fileLastGranule int64
		for i := len(pages) - 1; i >= 0; i-- {
			if pages[i].granulePos > 0 {
				fileLastGranule = pages[i].granulePos
				break
			}
		}

		headersDone := isFirst
		firstIncluded := false

		for pi := 0; pi < len(pages); pi++ {
			page := pages[pi]

			if !headersDone {
				if page.headerType&0x02 != 0 || page.granulePos <= 0 {
					continue
				}
				headersDone = true
			}

			modified := oggPage{
				headerType:     page.headerType,
				granulePos:     page.granulePos,
				serialNumber:   serialNumber,
				sequenceNumber: sequenceNum,
				segmentTable:   page.segmentTable,
				data:           page.data,
			}

			if !isFirst {
				modified.headerType &^= 0x02
				if !firstIncluded {
					modified.headerType &^= 0x01
					firstIncluded = true
				}
			}

			if !isFirst && page.granulePos > 0 {
				modified.granulePos = page.granulePos + granuleOffset
			}

			if !isLast {
				modified.headerType &^= 0x04
			}

			result = append(result, modified.serialize()...)
			sequenceNum++
		}

		granuleOffset += fileLastGranule
	}

	return result, nil
}
