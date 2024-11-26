package icons

import "github.com/maddalax/htmgo/framework/h"

func TrashIcon() *h.Element {
	return h.Svg(
		h.Class("h-full w-full"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Path(
			h.Attribute("d", "M10 11V17"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Path(
			h.Attribute("d", "M14 11V17"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Path(
			h.Attribute("d", "M4 7H20"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Path(
			h.Attribute("d", "M6 7H12H18V18C18 19.6569 16.6569 21 15 21H9C7.34315 21 6 19.6569 6 18V7Z"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Path(
			h.Attribute("d", "M9 5C9 3.89543 9.89543 3 11 3H13C14.1046 3 15 3.89543 15 5V7H9V5Z"),
			h.Attribute("stroke", "#000000"),
			h.Attribute("stroke-width", "2"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
	)
}

func Question() *h.Element {
	return h.Svg(
		h.Class("h-full w-full"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("stroke", "#000000"),
		h.Attribute("stroke-width", "2"),
		h.Attribute("stroke-linecap", "round"),
		h.Attribute("stroke-linejoin", "round"),
		h.Tag(
			"circle",
			h.Attribute("cx", "12"),
			h.Attribute("cy", "12"),
			h.Attribute("r", "10"),
		),
		h.Path(
			h.Attribute("d", "M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"),
		),
		h.Tag(
			"line",
			h.Attribute("x1", "12"),
			h.Attribute("y1", "17"),
			h.Attribute("x2", "12.01"),
			h.Attribute("y2", "17"),
		),
	)
}

func CheckMark() *h.Element {
	return h.Svg(
		h.Class("h-full w-full"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Tag("g",
			h.Attribute("stroke-width", "0"),
		),
		h.Tag("g",
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Tag("g",
			h.Path(
				h.Attribute("d", "M22.7048 4.95406C22.3143 4.56353 21.6811 4.56353 21.2906 4.95406L8.72696 17.5177C8.33643 17.9082 7.70327 17.9082 7.31274 17.5177L2.714 12.919C2.32348 12.5284 1.69031 12.5284 1.29979 12.919C0.909266 13.3095 0.909265 13.9427 1.29979 14.3332L5.90392 18.9289C7.07575 20.0986 8.97367 20.0978 10.1445 18.9271L22.7048 6.36827C23.0953 5.97775 23.0953 5.34458 22.7048 4.95406Z"),
				h.Attribute("fill", "#0F0F0F"),
			),
		),
	)
}

func ChevronDown() *h.Element {
	return h.Svg(
		h.Class("h-4 w-4"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("stroke", "currentColor"),
		h.Attribute("stroke-width", "2"),
		h.Attribute("stroke-linecap", "round"),
		h.Attribute("stroke-linejoin", "round"),
		h.Class("h-4 w-4 opacity-50"),
		h.Path(
			h.Attribute("d", "m6 9 6 6 6-6"),
		),
	)
}
