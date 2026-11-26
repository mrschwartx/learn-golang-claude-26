module contoh.com/app

go 1.24

require (
	contoh.com/pustaka v1.5.0
	contoh.com/pustaka/v2 v2.0.0
)

replace contoh.com/pustaka => ../pustaka-v1

replace contoh.com/pustaka/v2 => ../pustaka-v2
