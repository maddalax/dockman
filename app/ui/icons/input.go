package icons

import "github.com/maddalax/htmgo/framework/h"

func SearchIcon() *h.Element {
	return h.Svg(
		h.Class("h-5 w-5"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Tag(
			"g",
			h.Attribute("stroke-width", "0"),
		),
		h.Tag(
			"g",
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Tag(
			"g",
			h.Path(
				h.Attribute("d", "M15.7955 15.8111L21 21M18 10.5C18 14.6421 14.6421 18 10.5 18C6.35786 18 3 14.6421 3 10.5C3 6.35786 6.35786 3 10.5 3C14.6421 3 18 6.35786 18 10.5Z"),
				h.Attribute("stroke", "#000000"),
				h.Attribute("stroke-width", "2"),
				h.Attribute("stroke-linecap", "round"),
				h.Attribute("stroke-linejoin", "round"),
			),
		),
	)
}

func EyeIcon() *h.Element {
	return h.Svg(
		h.Class("h-4 h-4"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Tag(
			"g",
			h.Id("SVGRepo_bgCarrier"),
			h.Attribute("stroke-width", "0"),
		),
		h.Tag(
			"g",
			h.Id("SVGRepo_tracerCarrier"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Tag(
			"g",
			h.Id("SVGRepo_iconCarrier"),
			h.Path(
				h.Attribute("d", "M9 4.45962C9.91153 4.16968 10.9104 4 12 4C16.1819 4 19.028 6.49956 20.7251 8.70433C21.575 9.80853 22 10.3606 22 12C22 13.6394 21.575 14.1915 20.7251 15.2957C19.028 17.5004 16.1819 20 12 20C7.81811 20 4.97196 17.5004 3.27489 15.2957C2.42496 14.1915 2 13.6394 2 12C2 10.3606 2.42496 9.80853 3.27489 8.70433C3.75612 8.07914 4.32973 7.43025 5 6.82137"),
				h.Attribute("stroke", "#1C274C"),
				h.Attribute("stroke-width", "1.5"),
				h.Attribute("stroke-linecap", "round"),
			),
			h.Path(
				h.Attribute("d", "M15 12C15 13.6569 13.6569 15 12 15C10.3431 15 9 13.6569 9 12C9 10.3431 10.3431 9 12 9C13.6569 9 15 10.3431 15 12Z"),
				h.Attribute("stroke", "#1C274C"),
				h.Attribute("stroke-width", "1.5"),
			),
		),
	)
}
