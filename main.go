package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/gocolly/colly"
)

type Results struct {
	Results []struct {
		Title string `json:"titleNoFormatting"`
		URL   string `json:"unescapedUrl"`
	} `json:"results"`
}

func get_search_results(bird string) Results {
	bird = strings.Replace(bird, " ", "+", -1)

	client := &http.Client{}

	req, err := http.NewRequest("GET", string("https://cse.google.com/cse/element/v1?rsz=filtered_cse&num=10&hl=fr&source=gcsc&gss=.com&cselibv=8435450f13508ca1&cx=014496470795211077046%3AWMX431797713&q="+string(bird)+"&safe=off&cse_tok=AB-tC_7g-yFwYyx0deIBSf94gXCY%3A1709311385919&exp=cc&fexp=72497452&callback=google.search.cse.api7554"), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "cse.google.com")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cookie", "CONSENT=PENDING+416; __Secure-3PSID=g.a000gwhOmmMiKkPV3NOQLAJpomM5jL3Oc64OkvzdmItW-7gvQAqNEQCIpDsr3dHqM-CnMUHZowACgYKAe0SAQASFQHGX2MiZYC4-Bh9BpEbKvnGq0TGnhoVAUF8yKrvl4OLT3shcmyx0mTZUyAu0076; __Secure-3PAPISID=GaVrcj6Id6waauN4/AbAV0G5f6qNECecZT; NID=512=SFpNtzRXqBmXdDiUwqqnQEIDhLzpu8i4e3gcHMlwUPLycRjDipzIHCwEfQBe3svVV6J5FposuUKsCe6c_iHu81HpYUUSZfmRTGbRo_lCkQfml3px3ppYwSm4l3ZGGzpjLPQcK2a8osb83Hq87fD-rtvGlaMbOFzZ6ETUwL77PeLc9JZij4ZQk0HwIH6iKhIcVBjjKwUo8HX-rzdv0iyzrwKo2zheocUOBd3TPJMMjrCrym0GnVkbCI7USLMWwFKI5Pb11vqg1iqvRxW-RG_p5pbMLrF5Y8lJ8oevbEABts7pDnZPzI3wam_3ZJDVaW0mJxjTHMXJ6AYGO2K6Ylgjs5H4op3bF4YH8Mkvld6zLhIdB4Sg8yAenAe3-bkFG7NQK16StDcDpvcQxbcEu5RT02GjSxBd4o2sND2qIE-lpyxgzR7eOouZLX-ufiBJR_srbu-NadXwIihjzCQS0Ueu43I6_zPs8BU5NEMzmpUjccVW9N5J9NU_c0t6fWj4YXYN52Wy_jhpR6u4MtjZoiKKtpe6GIbz9XacLaiK-_lauGzPLv3eXvrdYHQlw38iOegF_oXSVQ5pqsZB69hGfce63VBC56zwGTw09Q3dr63zg9msKqIJr-j-Q1c0shGOFwPRvAn7EcTJR63mLA-BiKH7OQxiZ6ORW_ar9t_ANlbQZ9OJD6apQDgLFJYlG0YgPqaNMYgukjSsqqkZfZDcPAIXS_gzsu2nKI1lxMgi1QhQJhbbm2VJRvF1v_cQJhF-r0Dc4sLNDKYy4A; __Secure-3PSIDTS=sidts-CjEBYfD7Zwe8VPDpo8lMyXuYsC2oiG-IxflIbOgBmUP3pVchcU9M-vNw-q2hGefrHKcGEAA; __Secure-3PSIDCC=AKEyXzUYo_EheoIvP8UkVNAKv19-wpYozQo3IEJrJnl8XNEPbNlXkryxC-w-A1Skmr9fFrh3oA")
	req.Header.Set("referer", "https://www.oiseaux.net/")
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Opera GX";v="107", "Chromium";v="121"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "script")
	req.Header.Set("sec-fetch-mode", "no-cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 OPR/107.0.0.0")
	req.Header.Set("x-client-data", "CO3+ygE=")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyText)

	startIndex := strings.Index(bodyString, "{")
	jsonResponse := bodyString[startIndex:]
	jsonResponse = jsonResponse[:len(jsonResponse)-2]

	var results Results

	err = json.Unmarshal([]byte(jsonResponse), &results)
	if err != nil {
		panic(err)
	}

	return results
}

func get_bird_info(url string) []string {
	c := colly.NewCollector(colly.AllowedDomains("www.oiseaux.net"))

	var paragraphs []string
	c.OnHTML("div#description-esp", func(e *colly.HTMLElement) {
		p := e.ChildText("p")
		paragraphs = append(paragraphs, p)
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println("Erreur avec le site :", err)
	}
	return paragraphs
}

func main() {
	var search string
	huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Sur quel oiseau voulez vous avoir des informations ?").
				Value(&search).
				Validate(func(str string) error {
					if str == " " {
						return errors.New("veuillez renseigner un nom valide")
					}
					return nil
				}),
		),
	).Run()

	results := get_search_results(search)

	resultMap := make(map[string]string)
	for _, result := range results.Results {
		resultMap[result.Title] = result.URL
	}

	var numberRegex = regexp.MustCompile("[0-9]")

	title_list := []string{}
	for title := range resultMap {
		if strings.Contains(title, "-") && strings.Count(title, "-") == 1 && !strings.Contains(title, "Photos") && !strings.Contains(title, " : ") && !strings.Contains(title, ".net") {
			if !numberRegex.MatchString(title) {
				title_list = append(title_list, title)
			}
		}
	}

	if len(title_list) == 0 {
		panic("Aucun oiseau ayant ce nom n'a été trouvé")
	}

	var bird string
	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choisissez l'oiseau que vous cherchez").
				Options(huh.NewOptions(title_list...)...).
				Value(&bird),
		),
	).Run()

	desc := get_bird_info(resultMap[bird])

	huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(desc[0]),
		),
	).Run()
}
