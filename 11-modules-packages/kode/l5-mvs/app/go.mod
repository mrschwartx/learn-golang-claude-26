module contoh.com/app

go 1.24

require (
	contoh.com/deta v1.0.0
	contoh.com/detb v1.0.0
	contoh.com/lib v1.2.0
)

replace contoh.com/deta => ../deta

replace contoh.com/detb => ../detb

replace contoh.com/lib v1.1.0 => ../lib110

replace contoh.com/lib v1.2.0 => ../lib120

replace contoh.com/lib v1.3.0 => ../lib130
