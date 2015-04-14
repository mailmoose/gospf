package main

import (
	_ "fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetTerms(t *testing.T) {

	Convey("Testing getTerms()", t, func() {

		records := []struct {
			record     string
			directives []Directive
			modifiers  []Modifier
		}{
			{
				record:     "v=spf1 a -all",
				directives: []Directive{Directive{term: "a"}, Directive{term: "-all"}},
				modifiers:  []Modifier{},
			},
			{
				record:     "v=spf1 a:mail.example.com ~all",
				directives: []Directive{Directive{term: "a:mail.example.com"}, Directive{term: "~all"}},
				modifiers:  []Modifier{},
			},
			{
				record: "v=spf1 ip4:192.0.2.0/24 ip4:198.51.100.123 a -all",
				directives: []Directive{
					Directive{term: "ip4:192.0.2.0/24"},
					Directive{term: "ip4:198.51.100.123"},
					Directive{term: "a"},
					Directive{term: "-all"},
				},
				modifiers: []Modifier{},
			},
		}

		for _, record := range records {
			directives, modifiers, err := getTerms(record.record)
			So(err, ShouldEqual, nil)
			So(directives, ShouldResemble, record.directives)
			So(modifiers, ShouldResemble, record.modifiers)
		}

	})

}

func TestDirective(t *testing.T) {

	Convey("Testing directive.getQualifier()", t, func() {

		terms := []struct {
			d Directive
			q string
		}{
			{
				d: Directive{term: "ip4:192.0.2.0/24"},
				q: "",
			},
			{
				d: Directive{term: "-a"},
				q: "-",
			},
			{
				d: Directive{term: "+mx:mail.example.com"},
				q: "+",
			},
			{
				d: Directive{term: "~all"},
				q: "~",
			},
			{
				d: Directive{term: "?include:_spf.google.com"},
				q: "?",
			},
		}

		for _, term := range terms {
			So(term.d.getQualifier(), ShouldEqual, term.q)
		}

	})

	Convey("Testing directive.getMechanism()", t, func() {

		terms := []struct {
			d Directive
			m string
		}{
			{
				d: Directive{term: "ip4:192.0.2.0/24"},
				m: "ip4",
			},
			{
				d: Directive{term: "-a"},
				m: "a",
			},
			{
				d: Directive{term: "mx:mail.example.com"},
				m: "mx",
			},
			{
				d: Directive{term: "~all"},
				m: "all",
			},
			{
				d: Directive{term: "exists:example.com"},
				m: "exists",
			},
			{
				d: Directive{term: "include:_spf.google.com"},
				m: "include",
			},
		}

		for _, term := range terms {
			So(term.d.getMechanism(), ShouldEqual, term.m)
		}

	})

	Convey("Testing directive.getArguments()", t, func() {

		terms := []struct {
			d    Directive
			args []string
		}{
			{
				d:    Directive{term: "ip4:192.0.2.0/24"},
				args: []string{"24"},
			},
			{
				d:    Directive{term: "ip6:1080::8:800:68.0.3.1/96"},
				args: []string{"96"},
			},
			{
				d:    Directive{term: "a/32"},
				args: []string{"32"},
			},
			{
				d:    Directive{term: "a/24//96"},
				args: []string{"24", "96"},
			},
			{
				d:    Directive{term: "mx:foo.com//126"},
				args: []string{"", "126"},
			},
			{
				d:    Directive{term: "mx:foo.com/32"},
				args: []string{"32"},
			},
		}

		for _, term := range terms {
			So(term.d.getArguments(), ShouldResemble, term.args)
		}

	})

}
