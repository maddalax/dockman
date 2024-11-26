package icons

import (
	"github.com/maddalax/htmgo/framework/h"
	"strings"
)

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

func GitBranchIcon() *h.Element {
	return h.Svg(
		h.Class("h-4 w-4"),
		h.Attribute("viewBox", "0 0 24 24"),
		h.Attribute("fill", "none"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Tag("g",
			h.Id("SVGRepo_bgCarrier"),
			h.Attribute("stroke-width", "0"),
		),
		h.Tag("g",
			h.Id("SVGRepo_tracerCarrier"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Tag("g",
			h.Id("SVGRepo_iconCarrier"),
			h.Path(
				h.Attribute("fill-rule", "evenodd"),
				h.Attribute("clip-rule", "evenodd"),
				h.Attribute("d", "M6 5C6 4.44772 6.44772 4 7 4C7.55228 4 8 4.44772 8 5C8 5.55228 7.55228 6 7 6C6.44772 6 6 5.55228 6 5ZM8 7.82929C9.16519 7.41746 10 6.30622 10 5C10 3.34315 8.65685 2 7 2C5.34315 2 4 3.34315 4 5C4 6.30622 4.83481 7.41746 6 7.82929V16.1707C4.83481 16.5825 4 17.6938 4 19C4 20.6569 5.34315 22 7 22C8.65685 22 10 20.6569 10 19C10 17.7334 9.21506 16.6501 8.10508 16.2101C8.45179 14.9365 9.61653 14 11 14H13C16.3137 14 19 11.3137 19 8V7.82929C20.1652 7.41746 21 6.30622 21 5C21 3.34315 19.6569 2 18 2C16.3431 2 15 3.34315 15 5C15 6.30622 15.8348 7.41746 17 7.82929V8C17 10.2091 15.2091 12 13 12H11C9.87439 12 8.83566 12.3719 8 12.9996V7.82929ZM18 6C18.5523 6 19 5.55228 19 5C19 4.44772 18.5523 4 18 4C17.4477 4 17 4.44772 17 5C17 5.55228 17.4477 6 18 6ZM6 19C6 18.4477 6.44772 18 7 18C7.55228 18 8 18.4477 8 19C8 19.5523 7.55228 20 7 20C6.44772 20 6 19.5523 6 19Z"),
				h.Attribute("fill", "#000000"),
			),
		),
	)
}

func GithubIcon() *h.Element {
	return h.Svg(
		h.Class("h-4 w-4"),
		h.Attribute("viewBox", "0 0 20 20"),
		h.Attribute("version", "1.1"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Attribute("fill", "#000000"),
		h.Tag("g",
			h.Id("SVGRepo_bgCarrier"),
			h.Attribute("stroke-width", "0"),
		),
		h.Tag("g",
			h.Id("SVGRepo_tracerCarrier"),
			h.Attribute("stroke-linecap", "round"),
			h.Attribute("stroke-linejoin", "round"),
		),
		h.Tag("g",
			h.Tag("defs"),
			h.Tag("g",
				h.Attribute("stroke", "none"),
				h.Attribute("stroke-width", "1"),
				h.Attribute("fill", "none"),
				h.Attribute("fill-rule", "evenodd"),
				h.Tag("g",
					h.Attribute("transform", "translate(-140.000000, -7559.000000)"),
					h.Attribute("fill", "#000000"),
					h.Tag("g",
						h.Attribute("transform", "translate(56.000000, 160.000000)"),
						h.Path(
							h.Attribute("d", "M94,7399 C99.523,7399 104,7403.59 104,7409.253 C104,7413.782 101.138,7417.624 97.167,7418.981 C96.66,7419.082 96.48,7418.762 96.48,7418.489 C96.48,7418.151 96.492,7417.047 96.492,7415.675 C96.492,7414.719 96.172,7414.095 95.813,7413.777 C98.04,7413.523 100.38,7412.656 100.38,7408.718 C100.38,7407.598 99.992,7406.684 99.35,7405.966 C99.454,7405.707 99.797,7404.664 99.252,7403.252 C99.252,7403.252 98.414,7402.977 96.505,7404.303 C95.706,7404.076 94.85,7403.962 94,7403.958 C93.15,7403.962 92.295,7404.076 91.497,7404.303 C89.586,7402.977 88.746,7403.252 88.746,7403.252 C88.203,7404.664 88.546,7405.707 88.649,7405.966 C88.01,7406.684 87.619,7407.598 87.619,7408.718 C87.619,7412.646 89.954,7413.526 92.175,7413.785 C91.889,7414.041 91.63,7414.493 91.54,7415.156 C90.97,7415.418 89.522,7415.871 88.63,7414.304 C88.63,7414.304 88.101,7413.319 87.097,7413.247 C87.097,7413.247 86.122,7413.234 87.029,7413.87 C87.029,7413.87 87.684,7414.185 88.139,7415.37 C88.139,7415.37 88.726,7417.2 91.508,7416.58 C91.513,7417.437 91.522,7418.245 91.522,7418.489 C91.522,7418.76 91.338,7419.077 90.839,7418.982 C86.865,7417.627 84,7413.783 84,7409.253 C84,7403.59 88.478,7399 94,7399"),
						),
					),
				),
			),
		),
	)
}

func GitlabIcon() *h.Element {
	return h.Svg(
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
				h.Attribute("d", "M20.6866 13.8822L13.1333 19.1915C12.4433 19.6766 11.5231 19.6766 10.833 19.1915L3.27976 13.8822C3.16308 13.7992 3.07653 13.6826 3.0325 13.5492C2.98846 13.4157 2.9892 13.2722 3.03459 13.1392L4.03161 10.1514L6.02563 4.2154C6.04497 4.16703 6.07576 4.12372 6.11552 4.08893C6.18028 4.03172 6.26481 4 6.35252 4C6.44022 4 6.52475 4.03172 6.58951 4.08893C6.6315 4.12819 6.66244 4.17716 6.67941 4.2312L8.21482 8.78981C8.48886 9.60346 9.25164 10.1514 10.1102 10.1514H13.8549C14.7141 10.1514 15.4772 9.60271 15.7508 8.78828L17.2869 4.2154C17.3063 4.16703 17.3371 4.12372 17.3768 4.08893C17.4416 4.03172 17.5261 4 17.6138 4C17.7015 4 17.7861 4.03172 17.8508 4.08893C17.8928 4.12819 17.9238 4.17716 17.9407 4.2312L19.9347 10.1672L20.9726 13.1392C21.0139 13.2763 21.0084 13.4227 20.9569 13.5565C20.9053 13.6904 20.8105 13.8046 20.6866 13.8822Z"),
				h.Attribute("stroke", "#323232"),
				h.Attribute("stroke-width", "2"),
				h.Attribute("stroke-linecap", "round"),
				h.Attribute("stroke-linejoin", "round"),
			),
		),
	)
}

func CodeIcon() *h.Element {
	return h.Svg(
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
				h.Attribute("d", "M7 8L3 11.6923L7 16M17 8L21 11.6923L17 16M14 4L10 20"),
				h.Attribute("stroke", "#000000"),
				h.Attribute("stroke-width", "2"),
				h.Attribute("stroke-linecap", "round"),
				h.Attribute("stroke-linejoin", "round"),
			),
		),
	)
}

func GitProviderIcon(url string) *h.Element {
	if strings.Contains(url, "github.com") {
		return GithubIcon()
	}
	if strings.Contains(url, "gitlab.com") {
		return GitlabIcon()
	}
	return CodeIcon()
}
