// Copyright 2017, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsonfmt

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		in   string
		out  string
		err  error
		opts []Option
	}{{
		in:  "",
		out: "",
		err: jsonError{line: 1, column: 1, message: `unable to parse value: unexpected EOF`},
	}, {
		in:  `["]`,
		out: `["]`,
		err: jsonError{line: 1, column: 2, message: `unable to parse string: "\"]"`},
	}, {
		in:  "[\n\n\n]",
		out: `[]`,
	}, {
		in:  "{\n\n\n}",
		out: `{}`,
	}, {
		in: `{"firstName":"John","lastName":"Smith","isAlive":true,"age":25,"address":{"streetAddress":"21 2nd Street","city":"New York","state":"NY","postalCode":"10021-3100"},"phoneNumbers":[{"type":"home","number":"212 555-1234"},{"type":"office","number":"646 555-4567"},{"type":"mobile","number":"123 456-7890"}],"children":[],"spouse":null}`,
		out: `
			{
				"firstName": "John",
				"lastName":  "Smith",
				"isAlive":   true,
				"age":       25,
				"address": {
					"streetAddress": "21 2nd Street",
					"city":          "New York",
					"state":         "NY",
					"postalCode":    "10021-3100"
				},
				"phoneNumbers": [
					{"type": "home",   "number": "212 555-1234"},
					{"type": "office", "number": "646 555-4567"},
					{"type": "mobile", "number": "123 456-7890"}
				],
				"children": [],
				"spouse":   null
			}`,
	}, {
		in: `[[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}], [{"0123456789": "0123456789"}, {"0123456789": "0123456789"}], [{"0123456789": "0123456789"}, {"0123456789": "0123456789"}]]`,
		out: `
			[
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}],
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}],
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}]
			]`,
	}, {
		in: `[[{"0123456789012345678901234567890123456789": "0123456789"}, {"0123456789": "0123456789012345678901234567890123456789"}], [{"0123456789": "0123456789"}, {"0123456789": "0123456789"}], [{"0123456789": "0123456789"}, {"0123456789": "0123456789"}]]`,
		out: `
			[
				[
					{"0123456789012345678901234567890123456789": "0123456789"},
					{"0123456789": "0123456789012345678901234567890123456789"}
				],
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}],
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}]
			]`,
	}, {
		in: `
			{
				"Management": {
					"ServeAddress": "localhost:8080", "PasswordSalt": "", "PasswordHash": "",
					"SMTP": {"RelayServer": "mail.example.com:587", "Password":"abcdefghijklmnopqrstuvwxyz", "From":"noreply@example.com", "To":"noreply@example.com"},
				},

				// SSH comment.
				"SSH": {
					"KeyFiles":       ["key.priv"], // SSH key file
					"KnownHostFiles": [], // SSH known hosts file
				},

				"RateLimit":    "10Mi",
				"AutoSnapshot": {"Cron": "* * * * *", "Count": 3, "TimeZone": "Local"},
				"SendFlags":    ["-w"],
				"RecvFlags":    ["-s"],
				"Datasets": [{
					"AutoSnapshot": {"Cron": "0 6 * * *", "TimeZone": "Local", "Count": 30},
					"Source":  "//example.com/tank/fizz",
					"Mirrors": ["//foo.example.com/tank/replicas/fizz-drive"],
				}, {
					"Source":  "//example.com/tank/buzz",
					"Mirrors": ["//foo.example.com/tank/replicas/buzz-drive"],
				}, {
					// Seperate dataset so it has its own readonly setting
					"Source":  "//example.com/tank/users",
					"Mirrors": ["//foo.example.com/tank/replicas/users"],
				}],
			}`,
		out: `
			{
				"Management": {
					"ServeAddress": "localhost:8080", "PasswordSalt": "", "PasswordHash": "",
					"SMTP": {
						"RelayServer": "mail.example.com:587",
						"Password":    "abcdefghijklmnopqrstuvwxyz",
						"From":        "noreply@example.com",
						"To":          "noreply@example.com",
					},
				},

				// SSH comment.
				"SSH": {
					"KeyFiles":       ["key.priv"], // SSH key file
					"KnownHostFiles": [],           // SSH known hosts file
				},

				"RateLimit":    "10Mi",
				"AutoSnapshot": {"Cron": "* * * * *", "Count": 3, "TimeZone": "Local"},
				"SendFlags":    ["-w"],
				"RecvFlags":    ["-s"],
				"Datasets": [{
					"AutoSnapshot": {"Cron": "0 6 * * *", "TimeZone": "Local", "Count": 30},
					"Source":       "//example.com/tank/fizz",
					"Mirrors":      ["//foo.example.com/tank/replicas/fizz-drive"],
				}, {
					"Source":  "//example.com/tank/buzz",
					"Mirrors": ["//foo.example.com/tank/replicas/buzz-drive"],
				}, {
					// Seperate dataset so it has its own readonly setting
					"Source":  "//example.com/tank/users",
					"Mirrors": ["//foo.example.com/tank/replicas/users"],
				}],
			}`,
	}, {
		in: `
			{
				"Management": {
					"ServeAddress": "localhost:8080", "PasswordSalt": "", "PasswordHash": "",
					"SMTP": {"RelayServer": "mail.example.com:587", "Password":"abcdefghijklmnopqrstuvwxyz", "From":"noreply@example.com", "To":"noreply@example.com"},
				},

				// SSH comment.
				"SSH": {
					"KeyFiles":       ["key.priv"], // SSH key file
					"KnownHostFiles": [], // SSH known hosts file
				},

				"RateLimit":    "10Mi",
				"AutoSnapshot": {"Cron": "* * * * *", "Count": 3, "TimeZone": "Local"},
				"SendFlags":    ["-w"],
				"RecvFlags":    ["-s"],
				"Datasets": [{
					"AutoSnapshot": {"Cron": "0 6 * * *", "TimeZone": "Local", "Count": 30},
					"Source":  "//example.com/tank/fizz",
					"Mirrors": ["//foo.example.com/tank/replicas/fizz-drive"],
				}, {
					"Source":  "//example.com/tank/buzz",
					"Mirrors": ["//foo.example.com/tank/replicas/buzz-drive"],
				}, {
					// Seperate dataset so it has its own readonly setting
					"Source":  "//example.com/tank/users",
					"Mirrors": ["//foo.example.com/tank/replicas/users"],
				}],
			}`,
		out:  `{"Management":{"ServeAddress":"localhost:8080","PasswordSalt":"","PasswordHash":"","SMTP":{"RelayServer":"mail.example.com:587","Password":"abcdefghijklmnopqrstuvwxyz","From":"noreply@example.com","To":"noreply@example.com"}},"SSH":{"KeyFiles":["key.priv"],"KnownHostFiles":[]},"RateLimit":"10Mi","AutoSnapshot":{"Cron":"* * * * *","Count":3,"TimeZone":"Local"},"SendFlags":["-w"],"RecvFlags":["-s"],"Datasets":[{"AutoSnapshot":{"Cron":"0 6 * * *","TimeZone":"Local","Count":30},"Source":"//example.com/tank/fizz","Mirrors":["//foo.example.com/tank/replicas/fizz-drive"]},{"Source":"//example.com/tank/buzz","Mirrors":["//foo.example.com/tank/replicas/buzz-drive"]},{"Source":"//example.com/tank/users","Mirrors":["//foo.example.com/tank/replicas/users"]}]}`,
		opts: []Option{Minify()},
	}, {
		in: `{"Management":{"ServeAddress":"localhost:8080","PasswordSalt":"","PasswordHash":"","SMTP":{"RelayServer":"mail.example.com:587","Password":"abcdefghijklmnopqrstuvwxyz","From":"noreply@example.com","To":"noreply@example.com"}},"SSH":{"KeyFiles":["key.priv"],"KnownHostFiles":[]},"RateLimit":"10Mi","AutoSnapshot":{"Cron":"* * * * *","Count":3,"TimeZone":"Local"},"SendFlags":["-w"],"RecvFlags":["-s"],"Datasets":[{"AutoSnapshot":{"Cron":"0 6 * * *","TimeZone":"Local","Count":30},"Source":"//example.com/tank/fizz","Mirrors":["//foo.example.com/tank/replicas/fizz-drive"]},{"Source":"//example.com/tank/buzz","Mirrors":["//foo.example.com/tank/replicas/buzz-drive"]},{"Source":"//example.com/tank/users","Mirrors":["//foo.example.com/tank/replicas/users"]}]}`,
		out: `
			{
				"Management": {
					"ServeAddress": "localhost:8080",
					"PasswordSalt": "",
					"PasswordHash": "",
					"SMTP": {
						"RelayServer": "mail.example.com:587",
						"Password":    "abcdefghijklmnopqrstuvwxyz",
						"From":        "noreply@example.com",
						"To":          "noreply@example.com"
					}
				},
				"SSH":          {"KeyFiles": ["key.priv"], "KnownHostFiles": []},
				"RateLimit":    "10Mi",
				"AutoSnapshot": {"Cron": "* * * * *", "Count": 3, "TimeZone": "Local"},
				"SendFlags":    ["-w"],
				"RecvFlags":    ["-s"],
				"Datasets": [{
					"AutoSnapshot": {"Cron": "0 6 * * *", "TimeZone": "Local", "Count": 30},
					"Source":       "//example.com/tank/fizz",
					"Mirrors":      ["//foo.example.com/tank/replicas/fizz-drive"]
				}, {
					"Source":  "//example.com/tank/buzz",
					"Mirrors": ["//foo.example.com/tank/replicas/buzz-drive"]
				}, {
					"Source":  "//example.com/tank/users",
					"Mirrors": ["//foo.example.com/tank/replicas/users"]
				}]
			}`,
	}, {
		in: "[\n123456789,\n123456789,\n123456789,\n]",
		out: `
			[
				123456789,
				123456789,
				123456789,
			]`,
	}, {
		in: "[\n123456789,\n123456789,\n123456789,\n]",
		out: `
			[
				123456789,
				123456789,
				123456789
			]`,
		opts: []Option{Standardize()},
	}, {
		in:   "[\n123456789,\n123456789,\n123456789,\n]",
		out:  "[123456789,123456789,123456789]",
		opts: []Option{Minify()},
	}, {
		in: `{
			/* comment */
			"key"
			/* comment */
			:
			/* comment */
			"record"
			/* comment */
			,
			/* comment */
			"key"
			/* comment */
			:
			/* comment */
			"record"
			/* comment */
			,
		}`,
		out: `
			{
				/* comment */
				"key"
				/* comment */ :
					/* comment */
					"record"
					/* comment */ ,
				/* comment */
				"key"
				/* comment */ :
					/* comment */
					"record"
					/* comment */ ,
			}`,
	}, {
		in: `

					/*
					* Block comment.
					*/
					"Text"
		`,
		out: `
			/*
			 * Block comment.
			 */
			"Text"`,
	}, {
		in: `
			[
								{
									"fwafwa" /*ffawe*/:
							    		"fewafwaf",

					"fwafwafwae":




				                 		"fwafewa",},

			[/*comment*/
			{/*comment*/},
				{

				}




				],

					{"fwafwa":



							    		"fewafwaf",
					"fwafwafwae": "dwafewa",//fea
					"fwafwafwae"://fa
					"fwafewa",},

					{
						"fwafwa": 0.0000000000000000000033242000000,
					"fwafwafwae"


					:				"fwafewa",
					 },
					 ["fweafewa","faewfaew","afwefawe"/*
					 fewfaew
					 fewafewa*/]
			 				    ]`,
		out: `
			[
				{
					"fwafwa" /*ffawe*/ :
						"fewafwaf",

					"fwafwafwae":
						"fwafewa",
				},

				[ /*comment*/
					{ /*comment*/ },
					{},
				],

				{
					"fwafwa":
						"fewafwaf",
					"fwafwafwae": "dwafewa", //fea
					"fwafwafwae":            //fa
						"fwafewa",
				},

				{
					"fwafwa":     3.3242e-21,
					"fwafwafwae": "fwafewa",
				},
				["fweafewa", "faewfaew", "afwefawe", /*
				fewfaew
				fewafewa*/ ],
			]`,
	}, {
		in: `
			[
				{"keyX": [1,2,3,4,5]},
				{"keyXX": [1,2,3], "keyZ": {"subkey": "value"},},
				{"keyY": "val", "keyZZ": [[[[[[[1,2,3]]]]]]]},
			]`,
		out: `
			[
				{"keyX":  [1, 2, 3, 4, 5]},
				{"keyXX": [1, 2, 3], "keyZ":  {"subkey": "value"}},
				{"keyY":  "val",     "keyZZ": [[[[[[[1, 2, 3]]]]]]]},
			]`,
	}, {
		in: `
			{
				"key": "val01234567", // Comment 1
				"key01234567890123456789": "val0123456789", // Comment 2
				"key": "val", // Comment 3
				"key0123456789": "val0123", // Comment 4
			}`,
		out: `
			{
				"key":                     "val01234567",   // Comment 1
				"key01234567890123456789": "val0123456789", // Comment 2
				"key":                     "val",           // Comment 3
				"key0123456789":           "val0123",       // Comment 4
			}`,
	}, {
		in: `
			{
				"key012345678901234567890": "val0123456789", // Comment 2
				"key": "val01234567", // Comment 1
				"key": "val", // Comment 3
				"key0123456789": "val0123", // Comment 4
			}`,
		out: `
			{
				"key012345678901234567890": "val0123456789", // Comment 2
				"key":           "val01234567",              // Comment 1
				"key":           "val",     // Comment 3
				"key0123456789": "val0123", // Comment 4
			}`,
	}, {
		in:  `/**//**/{/**//**/"key"/**//**/:/**//**/"val"/**//**/}/**//**/`,
		out: `/**/ /**/ { /**/ /**/ "key" /**/ /**/ : /**/ /**/ "val" /**/ /**/ } /**/ /**/`,
	}, {
		in:  `{"PrimeNumbers": [{}, 2, 3, 5, 7, 11, 13, 17, 19, {}]}`,
		out: `{"PrimeNumbers": [{}, 2, 3, 5, 7, 11, 13, 17, 19, {}]}`,
	}, {
		in: `{"PrimeNumbers": [{}, 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, {}, 503, 509, 521, 523, 541, 547, 557, 563, 569, 571, 577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643, 647, 653, 659, 661, 673, 677, 683, 691, 701, 709, 719, 727, 733, 739, 743, 751, 757, 761, 769, 773, 787, 797, 809, 811, 821, 823, 827, 829, 839, 853, 857, 859, 863, 877, 881, 883, 887, 907, 911, 919, 929, 937, 941, 947, 953, 967, 971, 977, 983, 991, 997, {}]}`,
		out: `
			{"PrimeNumbers": [
				{},
				2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73,
				79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157,
				163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241,
				251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347,
				349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439,
				443, 449, 457, 461, 463, 467, 479, 487, 491, 499,
				{},
				503, 509, 521, 523, 541, 547, 557, 563, 569, 571, 577, 587, 593, 599, 601, 607,
				613, 617, 619, 631, 641, 643, 647, 653, 659, 661, 673, 677, 683, 691, 701, 709,
				719, 727, 733, 739, 743, 751, 757, 761, 769, 773, 787, 797, 809, 811, 821, 823,
				827, 829, 839, 853, 857, 859, 863, 877, 881, 883, 887, 907, 911, 919, 929, 937,
				941, 947, 953, 967, 971, 977, 983, 991, 997,
				{}
			]}`,
	}, {
		in: `
			{"PrimeNumbers": [
				// Group 1.
				{}, 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79,	83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223,
				227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383,
				389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463,
				467, 479, 487, 491, 499, {}, 503, 509, 521, 523, 541, 547, 557, 563, 569, 571, 577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643, 647,
				// Group 2.
				653, 659, 661, 673, 677, 683,
				691, 701, 709, 719, 727, 733, 739, 743, 751, 757, 761,
				769, 773, 787, 797, 809, 811, 821, 823, 827, 829, 839, 853, 857, 859, 863, 877, 881, 883,
				887, 907, 911, 919, 929, 937, 941, 947, 953, 967, 971, 977, 983, 991, 997, {}
			]}`,
		out: `
			{"PrimeNumbers": [
				// Group 1.
				{}, 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
				73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157,
				163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223,
				227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311,
				313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383,
				389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463,
				467, 479, 487, 491, 499, {}, 503, 509, 521, 523, 541, 547, 557, 563, 569, 571,
				577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643, 647,
				// Group 2.
				653, 659, 661, 673, 677, 683,
				691, 701, 709, 719, 727, 733, 739, 743, 751, 757, 761,
				769, 773, 787, 797, 809, 811, 821, 823, 827, 829, 839, 853, 857, 859, 863, 877,
				881, 883,
				887, 907, 911, 919, 929, 937, 941, 947, 953, 967, 971, 977, 983, 991, 997, {}
			]}`,
	}, {
		in: `
			{
				"key": "val", "key": "val", "key": "reallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallylongvalue",
			}`,
		out: `
			{
				"key": "val", "key": "val",
				"key": "reallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallylongvalue",
			}`,
	}}

	for i, tt := range tests {
		// Adjust output for leading tabs and newlines.
		want := strings.Join(strings.Split(tt.out, "\n\t\t\t"), "\n")
		if strings.HasPrefix(want, "\n") {
			want = want[1:] + "\n"
		}

		got, err := Format([]byte(tt.in), tt.opts...)
		if got := string(got); got != want || err != tt.err {
			diff := cmp.Diff(strings.Split(got, "\n"), strings.Split(want, "\n"))
			t.Errorf("test %d, Format output mismatch (-got +want):\n%s\ngot  `%v`\nwant `%v`", i, diff, got, want)
		}
		if err != tt.err {
			t.Errorf("test %d, Format error mismatch:\ngot  %v\nwant %v`", i, err, tt.err)

		}
	}
}
