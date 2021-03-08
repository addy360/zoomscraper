package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Phone holding values
type Phone struct {
	phoneName string
	phoneTag  string
	location  string
	price     string
}

func (p *Phone) sanitizeFields() {
	p.location = strings.TrimSpace(p.location)
	p.phoneName = strings.TrimSpace(p.phoneName)
	p.phoneTag = strings.TrimSpace(p.phoneTag)
	p.price = strings.TrimSpace(p.price)
}

type phones []Phone

type byPrice phones

func (byPrice byPrice) Len() int {
	return len(byPrice)
}

func extractPrice(price string) int {
	priceString := strings.ReplaceAll(price, ",", "")
	priceString = strings.Split(priceString, " ")[1]
	priceInt, _ := strconv.Atoi(priceString)
	return priceInt
}

func (byPrice byPrice) Less(i, j int) bool {
	return extractPrice(byPrice[i].price) < extractPrice(byPrice[j].price)
}

func (byPrice byPrice) Swap(i, j int) {
	byPrice[i], byPrice[j] = byPrice[j], byPrice[i]
}

func (p *phones) search(phonename string) {
	for _, v := range *p {
		if strings.Contains(strings.ToUpper(v.phoneName), strings.ToUpper(phonename)) {
			fmt.Printf("Phone name : %s\nPhone price : %s\nLocation : %s\n\n", v.phoneName, v.price, v.location)
		}
	}

}

func (p *phones) userSearch() {
	var query string
	print("\n\nphone  >>> ")
	fmt.Scanf("%s\n", &query)

	p.search(query)
}

func main() {
	print("\n\n MADE BY Addy360 \n\n> ")
	time.Sleep(time.Second)
	var scrappedPhones phones

	c := colly.NewCollector(
		colly.AllowedDomains("www.zoomtanzania.com"),
	)

	c.OnHTML("div.listings-cards__list-item", func(h *colly.HTMLElement) {
		phoneName := h.ChildText(".listing-card__header__title")
		phoneTag := h.ChildText(".listing-card__header__tags")
		location := h.ChildText(".listing-card__header__location")
		price := h.ChildText(".listing-card__price__value")
		phone := Phone{
			phoneName: phoneName,
			phoneTag:  phoneTag,
			location:  location,
			price:     price,
		}

		phone.sanitizeFields()
		scrappedPhones = append(scrappedPhones, phone)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println(string(r.Body))
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println(r.URL.String())
	})

	var pages int
	defaultPages := 10
	fmt.Printf("How many pages to search, Default is ( %d pages )  >>> ", defaultPages)
	fmt.Scanf("%d", &pages)
	if pages < 1 {

		pages = defaultPages
	}

	for i := 1; i < pages+1; i++ {
		if err := c.Visit(fmt.Sprintf("https://www.zoomtanzania.com/mobile-phones?p=%d", i)); err != nil {
			log.Fatal("Make sure you are connected to the internet")
		}
	}

	sort.Sort(byPrice(scrappedPhones))
	fmt.Printf("\n\nFOUND %d PHONES\n", len(scrappedPhones))
	for {
		scrappedPhones.userSearch()
	}

}
