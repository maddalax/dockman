package icons

import "github.com/maddalax/htmgo/framework/h"

func DockerIconBlack() *h.Element {
	return h.Svg(
		h.Class("w-full h-full"),
		h.Attribute("fill", "#000000"),
		h.Attribute("viewBox", "0 0 32 32"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Path(
			h.Attribute("d", "M6.427 23.031c-0.911 0-1.739-0.744-1.739-1.651s0.744-1.656 1.739-1.656c1 0 1.751 0.745 1.751 1.656 0 0.907-0.833 1.651-1.745 1.651zM27.776 14.016c-0.183-1.323-1-2.401-2.079-3.224l-0.421-0.333-0.339 0.411c-0.656 0.745-0.921 2.068-0.839 3.057 0.079 0.751 0.317 1.495 0.74 2.073-0.344 0.177-0.76 0.333-1.084 0.505-0.76 0.249-1.5 0.333-2.239 0.333h-21.385l-0.084 0.489c-0.156 1.579 0.084 3.229 0.751 4.724l0.328 0.579v0.077c2 3.313 5.557 4.803 9.437 4.803 7.459 0 13.573-3.224 16.473-10.177 1.901 0.083 3.819-0.412 4.719-2.235l0.24-0.411-0.396-0.251c-1.083-0.661-2.563-0.749-3.801-0.416l-0.027 0.005zM17.099 12.693h-3.239v3.228h3.239zM17.099 8.636h-3.239v3.228h3.239zM17.099 4.495h-3.239v3.229h3.239zM21.057 12.693h-3.219v3.228h3.229v-3.228zM9.063 12.693h-3.219v3.228h3.229v-3.228zM13.099 12.693h-3.197v3.228h3.219v-3.228zM5.063 12.693h-3.199v3.228h3.24v-3.228zM13.099 8.636h-3.197v3.228h3.219v-3.224zM9.041 8.636h-3.192v3.228h3.219v-3.224l-0.021-0.004z"),
		),
	)
}

func DockerFileIconBlack() *h.Element {
	return h.Svg(
		h.Class("h-full w-full"),
		h.Attribute("viewBox", "0 0 631 667"),
		h.Attribute("xmlns", "http://www.w3.org/2000/svg"),
		h.Attribute("style", "fill-rule:evenodd;clip-rule:evenodd;stroke-linecap:round;stroke-linejoin:round;"),
		h.Tag(
			"g",
			h.Path(
				h.Attribute("d", "M367.661,63.034L367.661,177.756C367.661,193.311 367.661,201.088 370.523,207.029C373.042,212.255 377.057,216.504 382,219.167C387.617,222.194 394.97,222.194 409.679,222.194L518.152,222.194M525.23,277.413L525.23,477.713C525.23,524.379 525.23,547.71 516.642,565.535C509.09,581.213 497.038,593.959 482.214,601.947C465.359,611.029 443.299,611.029 399.175,611.029L231.101,611.029C186.978,611.029 164.916,611.029 148.063,601.947C133.239,593.959 121.186,581.213 113.633,565.535C105.046,547.71 105.046,524.379 105.046,477.713L105.046,188.863C105.046,142.199 105.046,118.866 113.633,101.043C121.186,85.365 133.239,72.618 148.063,64.63C164.916,55.548 186.978,55.548 231.101,55.548L315.448,55.548C334.716,55.548 344.351,55.548 353.419,57.85C361.458,59.891 369.142,63.258 376.193,67.826C384.143,72.979 390.955,80.184 404.582,94.595L488.309,183.145C501.936,197.556 508.748,204.761 513.62,213.17C517.94,220.625 521.123,228.753 523.053,237.254C525.23,246.844 525.23,257.033 525.23,277.413Z"),
				h.Attribute("style", "fill:none;fill-rule:nonzero;stroke:black;stroke-width:54.06px;"),
			),
			h.Path(
				h.Attribute("d", "M219.228,477.984C209.736,477.984 201.108,470.232 201.108,460.781C201.108,451.331 208.86,443.526 219.228,443.526C229.648,443.526 237.473,451.289 237.473,460.781C237.473,470.232 228.793,477.984 219.29,477.984L219.228,477.984ZM441.677,384.051C439.771,370.266 431.258,359.033 420.015,350.458L415.628,346.988L412.096,351.271C405.261,359.033 402.499,372.818 403.354,383.124C404.177,390.949 406.657,398.701 411.064,404.724C407.48,406.568 403.146,408.193 399.77,409.985C391.851,412.58 384.14,413.455 376.44,413.455L153.615,413.455L152.74,418.55C151.114,435.003 153.615,452.196 160.565,467.773L163.983,473.806L163.983,474.608C184.822,509.129 221.885,524.654 262.313,524.654C340.034,524.654 403.739,491.061 433.956,418.613C453.764,419.478 473.749,414.32 483.127,395.325L485.628,391.042L481.501,388.427C470.217,381.54 454.796,380.623 441.896,384.093L441.615,384.145L441.677,384.051ZM330.427,370.266L296.677,370.266L296.677,403.9L330.427,403.9L330.427,370.266ZM330.427,327.993L296.677,327.993L296.677,361.628L330.427,361.628L330.427,327.993ZM330.427,284.845L296.677,284.845L296.677,318.49L330.427,318.49L330.427,284.845ZM371.772,370.266L338.127,370.266L338.127,403.9L371.772,403.9L371.772,370.266ZM246.798,370.266L213.153,370.266L213.153,403.9L246.798,403.9L246.798,370.266ZM288.977,370.266L255.436,370.266L255.436,403.9L288.977,403.9L288.977,370.266ZM205.443,370.266L171.683,370.266L171.683,403.9L205.443,403.9L205.443,370.266ZM288.748,327.993L255.436,327.993L255.436,361.628L288.977,361.628L288.977,328.035L288.748,327.993ZM246.527,327.993L213.205,327.993L213.205,361.628L246.746,361.628L246.746,328.035L246.527,327.993Z"),
				h.Attribute("style", "fill-rule:nonzero;transform: scale(0.8); transform-origin: center;"),
			),
		),
	)
}
